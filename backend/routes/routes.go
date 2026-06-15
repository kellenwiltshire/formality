package routes

import (
	"fmt"
	"formality/backend/forms"
	loadenv "formality/backend/load_env"
	sendmail "formality/backend/send_mail"
	smtp_settings "formality/backend/smtp"
	"formality/backend/submissions"
	"formality/backend/users"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() {
	r := mux.NewRouter()

	// User Routes
	r.HandleFunc("/users/{id}", users.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", users.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", users.DeleteUser).Methods("DELETE")

	r.HandleFunc("/users", users.GetAllUsers).Methods("GET")
	r.HandleFunc("/users", users.CreateUser).Methods("POST")

	// Form Routes
	r.HandleFunc("/forms/{id}", forms.GetForm).Methods("GET")
	r.HandleFunc("/forms/{id}", forms.UpdateForm).Methods("PUT")
	r.HandleFunc("/forms/{id}", forms.DeleteForm).Methods("DELETE")
	r.HandleFunc("/forms/{id}", submissions.CreateFormResponse).Methods("POST")

	r.HandleFunc("/forms", forms.GetAllFormsForUser).Methods("GET")
	r.HandleFunc("/forms", forms.CreateForm).Methods("POST")

	r.HandleFunc("/forms/{id}/responses", submissions.GetFormResponses).Methods("GET")
	r.HandleFunc("/forms/{form_id}/responses/{id}", submissions.GetFormResponse).Methods("GET")
	r.HandleFunc("/forms/{form_id}/responses/{id}", submissions.DeleteFormResponse).Methods("DELETE")

	// SMTP
	r.HandleFunc("/email-settings/{user_id}", smtp_settings.GetSMTPSettings).Methods("GET")
	r.HandleFunc("/email-settings/{user_id}", smtp_settings.CreateSMTPSettings).Methods("POST")
	r.HandleFunc("/email-settings/{user_id}", smtp_settings.UpdateSMTPSettings).Methods("PUT")
	r.HandleFunc("/email-settings/{user_id}", smtp_settings.DeleteSMTPSettings).Methods("DELETE")
	r.HandleFunc("/email-settings/{user_id}/test", sendmail.TestEmail).Methods("GET")


	// Start the Service
	port := loadenv.LoadDotEnvVariable("PORT")
	fmt.Println("The server is running on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}