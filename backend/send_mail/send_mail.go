package sendmail

import (
	"fmt"
	"formality/backend/database"
	"formality/backend/encrypt_text"
	"formality/backend/response"
	smtp_settings "formality/backend/smtp"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/gorilla/mux"
)

func PrepareEmail (user_id string, response_id string) {
	if user_id == "" || response_id == "" {
		//Return an error
		return
    }

	// Retrieve the SMTP settings, fail if not found
	smtp_settings := GetSMTPSetting(user_id)

	// Retrieve the Submission
	payload := GetFormResponse(response_id)

	// Send the Email
	success :=SendMail(smtp_settings, payload)

	// Update DB With status
	if(!success){
		// Update DB with failure
	} else {
		// Update DB with success
	}
}

func TestEmail (w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
    user_id := (params["user_id"])
	if user_id == "" {
        response.HttpResponse(w, "", 0, "Invalid ID", 400)
        return
    }

	// Retrieve the SMTP settings, fail if not found
	smtp_settings := GetSMTPSetting(user_id)

	payload := "This is a test email for Formality!"

	success := SendMail(smtp_settings, payload)

	if !success {
		response.HttpResponse(w, "", 0, "Unable to send test email", 500)
		return
	}
	response.HttpResponse(w, "", 1, "Test email sent", 200)
}

func GetSMTPSetting (user_id string) smtp_settings.SMTP_Settings {
	var smtp_settings smtp_settings.SMTP_Settings
	dbErr := database.Db.QueryRow("SELECT id, user_id, host, port, username, password_encrypted, encryption_type, updated_at, recipient_email, sender_email FROM smtp_settings WHERE user_id = $1", user_id).Scan(&smtp_settings.Id, &smtp_settings.User_id, &smtp_settings.Host, &smtp_settings.Port, &smtp_settings.Username, &smtp_settings.Password, &smtp_settings.Encryption_Type, &smtp_settings.Updated_At, &smtp_settings.Recipient_Email, &smtp_settings.Sender_Email)
	if dbErr != nil {
		fmt.Println(dbErr)
	}

	decrypted_pass, err := encrypt_text.DecryptAES(smtp_settings.Password)
	if err != nil {
		fmt.Println(err)
	}

	smtp_settings.Password = string(decrypted_pass)

	return smtp_settings
}

func GetFormResponse (response_id string) string {

	var submission string
	dbErr := database.Db.QueryRow("SELECT payload FROM form_submissions WHERE id = $2", response_id).Scan(&submission)
	if dbErr != nil {
		fmt.Println(dbErr)
	}
	return submission
}

func SendMail (smtp_settings smtp_settings.SMTP_Settings, payload string) bool {
	fmt.Println(smtp_settings.Recipient_Email)
	// Set up authentication information.
	auth := smtp.PlainAuth("", smtp_settings.Username, smtp_settings.Password, smtp_settings.Host)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{smtp_settings.Recipient_Email}
	msg := []byte("To: "+smtp_settings.Recipient_Email+"\r\n" +
		"Subject: New Form response from Formality!\r\n" +
		"\r\n" +
		payload,)

	fmt.Println(to)
	fmt.Println(msg)
	err := smtp.SendMail(smtp_settings.Host+":"+strconv.Itoa(smtp_settings.Port), auth, smtp_settings.Sender_Email, to, msg)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}