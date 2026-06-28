package listener

import (
	"fmt"
	"os"
	"time"

	"github.com/kellenwiltshire/formality/internal/app"

	"github.com/lib/pq"
)

func Listener(app *app.Application) error {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return fmt.Errorf("Must provide a Database Host")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		return fmt.Errorf("Must provide a Database Port")
	}
	user := os.Getenv("DB_USERNAME")
	if user == "" {
		return fmt.Errorf("Must provide a Database User")
	}
	dbname := os.Getenv("DB_DATABASE_NAME")
	if dbname == "" {
		return fmt.Errorf("Must provide a Database Name")
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		return fmt.Errorf("Must provide a Database Password")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		pass,
		host,
		port,
		dbname,
	)

	watcher := pq.NewListener(
		connStr,
		10*time.Second,
		time.Minute,
		nil,
	)

	err := watcher.Listen("form_submissions_inserts")
	if err != nil {
		return err
	}

	for n := range watcher.Notify {
		if n != nil {
			err := app.SendMailService.SendMail(n.Extra)
			if err != nil {
				app.Logger.Printf("ERROR: SendMailError %v", err)
			}
		}
	}
	return nil
}
