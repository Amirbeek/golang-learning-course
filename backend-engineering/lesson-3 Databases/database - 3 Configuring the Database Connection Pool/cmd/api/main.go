package main

import (
	"log"

	"github.com/amirbeek/social/internal/env"
	store2 "github.com/amirbeek/social/internal/store"
)

func main() {
	cfg := config{
		Addr: env.GetString("ADDR", ":8080"),
	}
	store := store2.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}
