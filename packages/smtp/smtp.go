package smtp_settings

import (
	"encoding/json"
	"fmt"
	"formality/packages/database"
	"formality/packages/encrypt_text"
	"formality/packages/response"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type SMTP_Settings struct {
	Id              int    `json:"id"`
	User_id         *int   `json:"user_id"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Encryption_Type string `json:"encryption_type"`
	Recipient_Email string `json:"recipient_email"`
	Sender_Email    string `json:"sender_email"`
	Updated_At      string `json:"updated_at"`
}

func GetSMTPSettings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user_id := params["user_id"]

	if user_id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	var smtp_settings SMTP_Settings
	dbErr := database.Db.QueryRow("SELECT id, user_id, host, port, username, encryption_type, recipient_email, sender_email, updated_at FROM smtp_settings WHERE user_id = $1", user_id).Scan(&smtp_settings.Id, &smtp_settings.User_id, &smtp_settings.Host, &smtp_settings.Port, &smtp_settings.Username, &smtp_settings.Encryption_Type, &smtp_settings.Recipient_Email, &smtp_settings.Sender_Email, &smtp_settings.Updated_At)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to get smtp settings for user", 500)
		return
	}
	response.HttpResponse(w, smtp_settings, 1, "", 200)

}

func CreateSMTPSettings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["user_id"]

	if id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	var smtp_settings SMTP_Settings
	_ = json.NewDecoder(r.Body).Decode(&smtp_settings)

	if smtp_settings.Host == "" || smtp_settings.Password == "" || smtp_settings.Port == 0 || smtp_settings.Username == "" || smtp_settings.Recipient_Email == "" || smtp_settings.Sender_Email == "" {
		response.HttpResponse(w, "", 0, "Missing parameter", 400)
		return
	}

	fmt.Println(smtp_settings.Password)

	encrypted_pass, err := encrypt_text.EncryptAES(smtp_settings.Password)
	if err != nil {
		response.HttpResponse(w, "", 0, "Encryption Error", 500)
		return
	}

	fmt.Println(encrypted_pass)

	_, err = database.Db.Exec("INSERT INTO smtp_settings (user_id, host, port, username, password_encrypted, encryption_type, recipient_email, sender_email) values ($1, $2, $3, $4, $5, $6, $7, $8)", id, smtp_settings.Host, smtp_settings.Port, smtp_settings.Username, encrypted_pass, smtp_settings.Encryption_Type, smtp_settings.Recipient_Email, smtp_settings.Sender_Email)
	if err != nil {
		fmt.Println(err)
		response.HttpResponse(w, "", 0, "Unable to create smtp settings", 500)
	} else {
		response.HttpResponse(w, "", 1, "smtp created", 200)
	}
}

func UpdateSMTPSettings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user_id := params["user_id"]

	if user_id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	var smtp_settings SMTP_Settings
	_ = json.NewDecoder(r.Body).Decode(&smtp_settings)

	if smtp_settings.Host == "" || smtp_settings.Password == "" || smtp_settings.Port == 0 || smtp_settings.Username == "" {
		response.HttpResponse(w, "", 0, "Missing parameter", 400)
		return
	}

	fmt.Println(smtp_settings.Password)

	encrypted_pass, err := encrypt_text.EncryptAES(smtp_settings.Password)
	if err != nil {
		response.HttpResponse(w, "", 0, "Encryption Error", 500)
		return
	}

	fmt.Println(encrypted_pass)

	time := time.Now()

	_, err = database.Db.Exec("UPDATE smtp_settings SET host = $1, port = $2, username = $3, password_encrypted = $4, encryption_type = $5, updated_at = $6, recipient_email = $7, sender_email = $8 WHERE user_id = $9", smtp_settings.Host, smtp_settings.Port, smtp_settings.Username, encrypted_pass, smtp_settings.Encryption_Type, time, smtp_settings.Recipient_Email, smtp_settings.Sender_Email, user_id)
	if err != nil {
		fmt.Println(err)
		response.HttpResponse(w, "", 0, "Unable to create smtp settings", 500)
	} else {
		response.HttpResponse(w, "", 1, "smtp created", 200)
	}

}

func DeleteSMTPSettings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user_id := (params["user_id"])
	if user_id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	_, dbErr := database.Db.Exec("DELETE FROM smtp_settings WHERE user_id = $1", user_id)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to delete smtp setting", 500)
		return
	}
	response.HttpResponse(w, "", 1, "SMTP Setting Deleted", 200)
}
