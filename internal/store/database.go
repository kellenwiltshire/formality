package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"strconv"

	"github.com/kellenwiltshire/formality/internal/util"

	"github.com/pressly/goose/v3"
)

func ConnectDatabase() (*sql.DB, error) {

	host, err := util.LoadDotEnvVariable("DB_HOST")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}
	envPort, err := util.LoadDotEnvVariable("DB_PORT")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}
	port, _ := strconv.Atoi(envPort)
	user, err := util.LoadDotEnvVariable("DB_USERNAME")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}
	dbname, err := util.LoadDotEnvVariable("DB_DATABASE_NAME")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}
	pass, err := util.LoadDotEnvVariable("DB_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	// set up postgres sql to open it.
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)

	db, err := sql.Open("postgres", psqlSetup)
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: open %w", err)
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
