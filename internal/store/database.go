package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"strconv"

	"github.com/pressly/goose/v3"
)

func ConnectDatabase() (*sql.DB, error) {

	host := os.Getenv("DB_HOST")
	if host == "" {
		return nil, fmt.Errorf("Must provide a Database Host")
	}
	port, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)
	if port == 0 || err != nil {
		return nil, fmt.Errorf("Must provide a Database Port")
	}
	user := os.Getenv("DB_USERNAME")
	if user == "" {
		return nil, fmt.Errorf("Must provide a Database User")
	}
	dbname := os.Getenv("DB_DATABASE_NAME")
	if dbname == "" {
		return nil, fmt.Errorf("Must provide a Database Name")
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		return nil, fmt.Errorf("Must provide a Database Password")
	}

	// set up postgres sql to open it.
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)

	db, err := sql.Open("postgres", psqlSetup)
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: ping %w", err)
	}

	fmt.Println("Connected to Database...")

	return db, nil
}

func MigrateFs(db *sql.DB, migrationsFs fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFs)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
