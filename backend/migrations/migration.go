package migration

import (
	"database/sql"
	"fmt"
	loadenv "formality/backend/load_env"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Migration() {
	// Load DB Env Variables
	dbUser := loadenv.LoadDotEnvVariable("DB_USERNAME")
	dbPass := loadenv.LoadDotEnvVariable("DB_PASSWORD")
	dbHost := loadenv.LoadDotEnvVariable("DB_HOST")
	dbPort := loadenv.LoadDotEnvVariable("DB_PORT")
	dbName := loadenv.LoadDotEnvVariable("DB_DATABASE_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run migrations automatically
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations completed successfully!")
}
