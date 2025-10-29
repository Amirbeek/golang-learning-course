package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amirbeek/social/docs"
	"github.com/amirbeek/social/internal/auth"
	"github.com/amirbeek/social/internal/mailer"
	"github.com/amirbeek/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

type mailConfig struct {
	sendGrid  sendGridConfig
	fromEmail string
	exp       time.Duration
	mailTrap  mailTrapConfig
}

type sendGridConfig struct {
	apiKey string
}
type mailTrapConfig struct {
	apiKey string
}

type basicConfig struct {
	user string
	pass string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type config struct {
	Addr        string
	DB          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
}

type dbConfig struct {
	Addr         string // The database connection address (DSN or host:port)
	MaxOpenConns int    // The maximum number of open connections allowed to the database at one time
	MaxIdleConns int    // The maximum number of idle (unused) connections that can remain open in the pool
	MaxIdleTime  string // The maximum amount of time a connection can remain idle before being closed (e.g., "5m" = 5 minutes)
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running, and ready to accept connections"))
	})

	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware).Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.Addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				r.Delete("/", app.deleteUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)

				//r.Patch("/", app.updatePostHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/feed", app.getUserFeedHandler)
			})
		})
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})

	})
	return r
}

func (app *application) run(mux *chi.Mux) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	if app.config.Addr == "" {
		log.Println("Warning: Addr is empty, defaulting to :8081")
		app.config.Addr = ":8081"
	}

	srv := http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	app.logger.Infow("Server has started", "addr", app.config.Addr, "env", app.config.env)

	return srv.ListenAndServe()
}
