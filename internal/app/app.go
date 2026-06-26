package app

import (
	"database/sql"
	"formality/internal/api"
	"formality/internal/middleware"
	"formality/internal/service"
	"formality/internal/store"
	"formality/migrations"
	"log"
	"os"
)

type Application struct {
	Logger             *log.Logger
	Db                 *sql.DB
	UserHandler        *api.UserHandler
	FormHandler        *api.FormHandler
	SMTPHandler        *api.SmtpHandler
	SubmissionsHandler *api.SubmissionHandler
	TokenHandler       *api.TokenHandler

	SendMailService *service.SendMailService

	Middleware middleware.UserMiddleware
}

func NewApplication() (*Application, error) {
	pgDb, err := store.ConnectDatabase()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFs(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	userStore := store.NewPostgresUserStore(pgDb)
	formStore := store.NewPostgresFormStore(pgDb)
	smtpStore := store.NewPostgresSmtpStore(pgDb)
	submissionsStore := store.NewPostgresSubmissionsStore(pgDb)
	tokenStore := store.NewPostgresTokenStore(pgDb)

	userHandler := api.NewUserHandler(userStore, logger)
	formHandler := api.NewFormHandler(formStore, logger)
	submissionsHandler := api.NewSubmissionHandler(submissionsStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	sendMailService := service.NewSendMailService(formStore, submissionsStore, smtpStore, logger)
	smtpHandler := api.NewSmtpHandler(smtpStore, *sendMailService, logger)

	middleware := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger:             logger,
		Db:                 pgDb,
		UserHandler:        userHandler,
		FormHandler:        formHandler,
		SMTPHandler:        smtpHandler,
		SubmissionsHandler: submissionsHandler,
		TokenHandler:       tokenHandler,

		SendMailService: sendMailService,

		Middleware: middleware,
	}

	return app, nil
}
