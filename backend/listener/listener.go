package listener

import (
	"fmt"
	loadenv "formality/backend/load_env"
	sendmail "formality/backend/send_mail"
	"time"

	"github.com/lib/pq"
)

func Listener() {
	host := loadenv.LoadDotEnvVariable("DB_HOST")
	port := loadenv.LoadDotEnvVariable("DB_PORT")
	user := loadenv.LoadDotEnvVariable("DB_USERNAME")
	dbname := loadenv.LoadDotEnvVariable("DB_DATABASE_NAME")
	pass := loadenv.LoadDotEnvVariable("DB_PASSWORD")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		pass,
		host,
		port,
		dbname,
	)

	fmt.Println(connStr)

	watcher := pq.NewListener(
		connStr,
		10*time.Second,
		time.Minute,
		nil,
	)

	err := watcher.Listen("form_submissions_inserts")
	if err != nil {
		fmt.Println(err)
	}

	for n := range watcher.Notify {
		if n != nil {
			fmt.Println(n.Extra)
			sendmail.PrepareEmail(n.Extra)
		}
	}
}