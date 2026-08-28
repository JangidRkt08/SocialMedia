package main

import (
	"log"

	"github.com/JangidRkt08/SocialMedia/internal/db"
	"github.com/JangidRkt08/SocialMedia/internal/env"
	"github.com/JangidRkt08/SocialMedia/internal/store"
)

const version = "0.0.2"

//	@title			GOSocial API
//	@description	API for GOSocial, xkcnkvxa social network for Gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host						localhost:3000
// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorisation
// @description
func main() {
	cfg := config{
		addr:   env.GetString("APP_ADDR", ":3001"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:3000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15min"),
		},
		env: env.GetString("ENV", "development"),
	}
	db, err := db.New(cfg.db.addr,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("db connected")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
