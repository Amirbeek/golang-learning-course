package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	Addr string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)    // logs each request (method, path, time, etc.).
	r.Use(middleware.Recoverer) // catches panics so the server doesn’t crash.
	r.Use(middleware.RequestID) // adds a unique ID to each request (useful for tracking).
	r.Use(middleware.RealIP)    // gets the real client IP address (even behind proxies).

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome World"))
	})
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})

	return r
}

func (app *application) run(mux *chi.Mux) error {

	cfg := http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: 15 * time.Second, // max time to write a response.
		ReadTimeout:  15 * time.Second, // max time to read a request
		IdleTimeout:  60 * time.Second, // max time to keep a connection open when idle.
	}

	log.Printf("server has started at %s", app.config.Addr)
	return cfg.ListenAndServe()
}
