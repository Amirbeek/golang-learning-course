package main

import (
	"log"

	"github.com/amirbeek/social/internal/env"
)

func main() {
	cfg := config{
		Addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}
