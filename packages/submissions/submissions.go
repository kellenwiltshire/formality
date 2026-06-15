package submissions

import (
	"encoding/json"
	"fmt"
	"formality/packages/database"
	"formality/packages/response"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Submission struct {
	Id           int    `json:"id"`
	FormId       string `json:"form_id"`
	Payload      string `json:"payload"`
	Submitted_at string `json:"submitted_at"`
	Status       string `json:"status"`
}

func GetFormResponses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	var submissions []Submission
	rows, dbErr := database.Db.Query("SELECT id, form_id, payload, submitted_at, status FROM form_submissions WHERE form_id = $1", id)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to get all submissions", 500)
		return
	}
	for rows.Next() {
		var submission Submission
		if err := rows.Scan(&submission.Id, &submission.FormId, &submission.Payload, &submission.Submitted_at, &submission.Status); err != nil {
			fmt.Println(err)
			response.HttpResponse(w, "", 0, "Unable to get all form rows", 500)
			return
		}
		submissions = append(submissions, submission)
	}
	response.HttpResponse(w, submissions, 1, "", 200)

}

func GetFormResponse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	form_id := params["form_id"]

	if id == "" || form_id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	var submission Submission
	dbErr := database.Db.QueryRow("SELECT id, form_id, payload, submitted_at, status FROM form_submissions WHERE form_id = $1 AND id = $2", form_id, id).Scan(&submission.Id, &submission.FormId, &submission.Payload, &submission.Submitted_at, &submission.Status)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to get all submissions", 500)
		return
	}
	response.HttpResponse(w, submission, 1, "", 200)

}

func CreateFormResponse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	// Read the raw HTTP body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate that it is actually valid JSON
	if !json.Valid(body) {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	var response_id int64
	error := database.Db.QueryRow("INSERT INTO form_submissions (form_id, payload) values ($1, $2) RETURNING id", id, body).Scan(&response_id)
	if error != nil {
		fmt.Println(error)
		response.HttpResponse(w, "", 0, "Unable to create new submission", 500)
		return
	} else {
		response.HttpResponse(w, "", 1, "submission created", 200)
	}

}

func DeleteFormResponse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	form_id := params["form_id"]

	if id == "" || form_id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	_, dbErr := database.Db.Exec("DELETE FROM form_submissions WHERE id = $1 AND form_id = $2", id, form_id)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to delete submission", 500)
		return
	}
	response.HttpResponse(w, "", 1, "Submission Deleted", 200)

}
