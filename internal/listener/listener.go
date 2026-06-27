package listener

import (
	"fmt"
	"formality/internal/app"
	"formality/internal/util"
	"time"

	"github.com/lib/pq"
)

func Listener(app *app.Application) error {
	host, err := util.LoadDotEnvVariable("DB_HOST")
	if err != nil {
		return err
	}
	port, err := util.LoadDotEnvVariable("DB_PORT")
	if err != nil {
		return err
	}
	user, err := util.LoadDotEnvVariable("DB_USERNAME")
	if err != nil {
		return err
	}
	dbname, err := util.LoadDotEnvVariable("DB_DATABASE_NAME")
	if err != nil {
		return err
	}
	pass, err := util.LoadDotEnvVariable("DB_PASSWORD")
	if err != nil {
		return err
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

	err = watcher.Listen("form_submissions_inserts")
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
