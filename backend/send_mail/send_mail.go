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

func PrepareEmail (response_id string) {
	if  response_id == "" {
		error := fmt.Errorf("Error: No response id provided")
		updateDBError(response_id, error)
    }

	user_id, idErr := GetUserId(response_id)
	if idErr != nil {
		updateDBError(response_id, idErr)
	}


	// Retrieve the SMTP settings, fail if not found
	smtp_settings, smtpErr := GetSMTPSetting(user_id)
	if smtpErr != nil {
		updateDBError(response_id, smtpErr)
	}

	// Retrieve the Submission
	payload, form_id, payloadErr := GetFormResponse(response_id)
	if payloadErr != nil {
		updateDBError(response_id, payloadErr)
	}

	// Retrieve the email attached to the form
	email, emailErr := GetFormEmail(form_id)
	if emailErr != nil {
		updateDBError(response_id, emailErr)
	}

	smtp_settings.Recipient_Email = email //re-assign the recipient email to match the form specified one

	// Send the Email
	success, sendErr := SendMail(smtp_settings, payload)
	if sendErr != nil {
		updateDBError(response_id, sendErr)
	}
	

	// Update DB With status
	if success {
		_, dbErr := database.Db.Exec("UPDATE form_submissions SET status = $1 WHERE id = $2", "dispatched", response_id)
		if dbErr != nil {
			fmt.Println(dbErr)
		}
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
	smtp_settings, err := GetSMTPSetting(user_id)
	if err != nil {
		response.HttpResponse(w, "", 0, "Unable to send test email", 500)
		fmt.Println(err)
		return
	}

	payload := "This is a test email for Formality!"

	_, sendErr := SendMail(smtp_settings, payload)

	if sendErr != nil {
		response.HttpResponse(w, "", 0, "Unable to send test email", 500)
		return
	}
	response.HttpResponse(w, "", 1, "Test email sent", 200)
}

func GetUserId (response_id string) (string, error) {
	var form_id string
	dbErr := database.Db.QueryRow("SELECT form_id FROM form_submissions WHERE id = $1", response_id).Scan(&form_id)
	if dbErr != nil {
		fmt.Println(dbErr)
		return "", dbErr
	}

	var user_id int
	err := database.Db.QueryRow("SELECT user_id FROM forms WHERE id = $1", form_id).Scan(&user_id)
	if err != nil {
		fmt.Println(err)
		return "", dbErr
	}

	return strconv.Itoa(user_id), nil

}

func GetSMTPSetting (user_id string) (smtp_settings.SMTP_Settings, error) {
	var smtp_settings smtp_settings.SMTP_Settings
	dbErr := database.Db.QueryRow("SELECT id, user_id, host, port, username, password_encrypted, encryption_type, updated_at, recipient_email, sender_email FROM smtp_settings WHERE user_id = $1", user_id).Scan(&smtp_settings.Id, &smtp_settings.User_id, &smtp_settings.Host, &smtp_settings.Port, &smtp_settings.Username, &smtp_settings.Password, &smtp_settings.Encryption_Type, &smtp_settings.Updated_At, &smtp_settings.Recipient_Email, &smtp_settings.Sender_Email)
	if dbErr != nil {
		fmt.Println(dbErr)
		return smtp_settings, dbErr
	}

	decrypted_pass, err := encrypt_text.DecryptAES(smtp_settings.Password)
	if err != nil {
		fmt.Println(err)
		return smtp_settings, dbErr
	}

	smtp_settings.Password = string(decrypted_pass)

	return smtp_settings, nil
}

func GetFormResponse (response_id string) (string, string, error) {

	var payload string
	var form_id int
	dbErr := database.Db.QueryRow("SELECT payload, form_id FROM form_submissions WHERE id = $1", response_id).Scan(&payload, &form_id)
	if dbErr != nil {
		fmt.Println(dbErr)
		return "", "", dbErr
	}

	return payload, strconv.Itoa(form_id), nil
}

func GetFormEmail (form_id string) (string, error) {

	var email string
	dbErr := database.Db.QueryRow("SELECT target_email FROM forms WHERE id = $1", form_id).Scan(&email)
	if dbErr != nil {
		fmt.Println(dbErr)
		return "", dbErr
	}

	return email, nil
}

func SendMail (smtp_settings smtp_settings.SMTP_Settings, payload string) (bool, error) {
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
		return false, err
	}

	return true, nil
}

func updateDBError (response_id string, error error) {
	fmt.Println("Error:", error.Error())
	_, dbErr := database.Db.Exec("UPDATE form_submissions SET status = $1 WHERE id = $2", "error", response_id)
	if dbErr != nil {
		fmt.Println(dbErr)
	}
}