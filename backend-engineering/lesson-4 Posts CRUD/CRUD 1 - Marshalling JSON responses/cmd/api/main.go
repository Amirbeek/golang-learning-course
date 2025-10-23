package main

import (
	"log"

	"github.com/amirbeek/social/internal/db"
	"github.com/amirbeek/social/internal/env"
	store2 "github.com/amirbeek/social/internal/store"
)

const version string = "v0.1"

func main() {
	// Load configuration from environment
	cfg := config{
		Addr: env.GetString("ADDR", ":8081"), // default port :8081
		DB: dbConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://supervillager:adminpassword@localhost:5433/social?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 5),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Connect to database
	dbConn, err := db.New(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxIdleTime)
	if err != nil {
		log.Panicf("Database connection failed: %v", err)
	}
	defer dbConn.Close()
	log.Println("Database connection pool established")

	//  storage layer
	store := store2.NewStorage(dbConn)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Printf("Server starting at http://localhost%s...", cfg.Addr)
	if err := app.run(mux); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
