package main

import (
	"log"
	"net/http"
	"time"

	"github.com/amirbeek/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	Addr string
	DB   dbConfig
	env  string
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

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome World"))
		})
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
		})
	})

	return r
}

func (app *application) run(mux *chi.Mux) error {
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

	log.Printf("Server has started at http://localhost%s", app.config.Addr)
	return srv.ListenAndServe()
}
