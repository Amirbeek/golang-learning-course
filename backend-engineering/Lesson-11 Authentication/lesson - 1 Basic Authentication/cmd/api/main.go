package main

import (
	"time"

	"github.com/amirbeek/social/internal/db"
	"github.com/amirbeek/social/internal/env"
	mailer "github.com/amirbeek/social/internal/mailer"
	store2 "github.com/amirbeek/social/internal/store"
	"go.uber.org/zap"
)

const version string = "v0.1"

//	@title			GopherSocial
//	@description	API for GopherSocial, a social network for gahpers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8081
//	@BasePath	/v1

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	// Load configuration from environment
	cfg := config{
		Addr:        env.GetString("ADDR", ":8081"), // default port :8081
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8081"),
		frontendURL: env.GetString("FRONTEND_URL", "localhost:8081"),
		DB: dbConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://supervillager:adminpassword@localhost:5433/social?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 5),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", "hello@amirbekshomurodov.me"),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAIL_TRAP_API_KEY", "e9ae7e7015894ca627fb0a83ce47da15"),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Connect to database
	dbConn, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime)
	if err != nil {
		logger.Fatal("Database connection failed: %v", err)
	}
	defer dbConn.Close()
	logger.Info("Database connection pool established")

	//  storage layer
	store := store2.NewStorage(dbConn)

	//mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal("Mailtrap client failed: %v", err)
	}

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailtrap,
	}

	mux := app.mount()

	logger.Info("Server starting at http://localhost%s...", cfg.Addr)

	if err := app.run(mux); err != nil {
		logger.Fatalf("Server stopped with error: %v", err)
	}
}
