package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func LoadDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	return os.Getenv(key)
}

func Migration() {
	// Load DB Env Variables
	dbUser := LoadDotEnvVariable("DB_USERNAME")
	dbPass := LoadDotEnvVariable("DB_PASSWORD")
	dbHost := LoadDotEnvVariable("DB_HOST")
	dbPort := LoadDotEnvVariable("DB_PORT")
	dbName := LoadDotEnvVariable("DB_DATABASE_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	fmt.Println(connStr)

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
