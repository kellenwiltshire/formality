package routes

import (
	"formality/internal/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		// User Routes
		r.Get("/api/user", app.Middleware.RequireUser(app.UserHandler.HandleGetUser))
		r.Put("/api/user", app.Middleware.RequireUser(app.UserHandler.HandleUpdateUser))
		r.Delete("/api/user{id}", app.Middleware.RequireUser(app.UserHandler.HandleDeleteUser))

		// Form Routes
		r.Get("/api/forms/{form_id}", app.Middleware.RequireUser(app.FormHandler.HandleGetForm))
		r.Put("/api/forms/{form_id}", app.Middleware.RequireUser(app.FormHandler.HandleUpdateForm))
		r.Delete("/api/forms/{form_id}", app.Middleware.RequireUser(app.FormHandler.HandleDeleteForm))

		r.Get("/api/forms", app.Middleware.RequireUser(app.FormHandler.HandleGetAllFormsForUser))
		r.Post("/api/forms", app.Middleware.RequireUser(app.FormHandler.HandleCreateForm))

		r.Get("/api/forms/{form_id}/responses", app.Middleware.RequireUser(app.SubmissionsHandler.HandleGetFormSubmissions))
		r.Get("/api/forms/{form_id}/responses/{submission_id}", app.Middleware.RequireUser(app.SubmissionsHandler.HandleGetFormSubmissionById))
		r.Delete("/api/forms/{form_id}/responses/{submission_id}", app.Middleware.RequireUser(app.SubmissionsHandler.HandleDeleteFormSubmission))

		// SMTP
		r.Get("/api/email-settings/", app.Middleware.RequireUser(app.SMTPHandler.HandleGetSMTPSettings))
		r.Post("/api/email-settings/", app.Middleware.RequireUser(app.SMTPHandler.HandleCreateSmtpSettings))
		r.Put("/api/email-settings/", app.Middleware.RequireUser(app.SMTPHandler.HandleUpdateSmtpSettings))
		r.Delete("/api/email-settings/", app.Middleware.RequireUser(app.SMTPHandler.HandleDeleteSmtpSetting))
		r.Get("/api/email-settings/test", app.Middleware.RequireUser(app.SMTPHandler.HandleTestEmail))

	})

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.AuthenticateAdmin)

		// Admin User Routes
		r.Get("/api/admin/getUsers", app.Middleware.RequireAdmin(app.UserHandler.HandleGetAllUsers))
		r.Post("/api/admin/createUser", app.Middleware.RequireAdmin(app.UserHandler.HandleCreateUser))
		r.Get("/api/deleteUser/{id}", app.Middleware.RequireAdmin(app.UserHandler.HandleDeleteUser))

	})

	r.Post("/api/forms/{form_id}", app.SubmissionsHandler.HandleCreateSubmission)

	// Login
	r.Post("/api/auth/login", app.TokenHandler.HandleCreateToken)
	r.Get("/api/auth/logout", app.TokenHandler.HandleDeleteTokens)

	return r
}
