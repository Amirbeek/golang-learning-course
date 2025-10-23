package main

import (
	"log"

	"github.com/amirbeek/social/internal/db"
	"github.com/amirbeek/social/internal/env"
	store2 "github.com/amirbeek/social/internal/store"
)

func main() {
	cfg := config{
		Addr: env.GetString("ADDR", ":5433"),
		DB: dbConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://supervillager:adminpassword@localhost:5433/social?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 5),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	defer db.Close() // always close db
	log.Println("Database connection pool established")

	if err != nil {
		log.Panic(err)
	}

	store := store2.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
