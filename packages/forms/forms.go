package forms

import (
	"encoding/json"
	"fmt"
	"formality/packages/database"
	"formality/packages/response"
	"net/http"

	"github.com/gorilla/mux"
)

type Form struct {
	Id           string `json:"id"`
	User_id      int    `json:"user_id"`
	Name         string `json:"name"`
	Target_email string `json:"email"`
	Created_at   string `json:"created"`
}

func GetAllFormsForUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get("id")
	if id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	var forms []Form
	rows, dbErr := database.Db.Query("SELECT id, user_id, name, target_email, created_at FROM forms WHERE user_id = $1", id)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to get all forms", 500)
	}
	for rows.Next() {
		var form Form
		if err := rows.Scan(&form.Id, &form.User_id, &form.Name, &form.Target_email, &form.Created_at); err != nil {
			fmt.Println(err)
			response.HttpResponse(w, "", 0, "Unable to get all form rows", 500)
			return
		}
		forms = append(forms, form)
	}
	response.HttpResponse(w, forms, 1, "", 200)
}

func CreateForm(w http.ResponseWriter, r *http.Request) {
	var newForm Form
	_ = json.NewDecoder(r.Body).Decode(&newForm)

	if newForm.Name == "" || newForm.Target_email == "" {
		response.HttpResponse(w, "", 0, "Please provide a name and email", 400)
		return
	}

	_, err := database.Db.Exec("INSERT INTO forms (user_id, name, target_email) values ($1, $2, $3)", newForm.User_id, newForm.Name, newForm.Target_email)
	if err != nil {
		fmt.Println(err)
		// We try again since the failure could be due to a collision on the created form id
		// Not the best solution, but works for now
		_, err := database.Db.Exec("INSERT INTO forms (user_id, name, target_email) values ($1, $2, $3)", newForm.User_id, newForm.Name, newForm.Target_email)
		if err != nil {
			fmt.Println(err)
			response.HttpResponse(w, "", 0, "Unable to create new form", 500)
			return
		} else {
			response.HttpResponse(w, "", 1, "Form created", 200)
			return
		}
	} else {
		response.HttpResponse(w, "", 1, "Form created", 200)
	}
}

func GetForm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	var form Form
	dbErr := database.Db.QueryRow("SELECT id, user_id, name, target_email, created_at FROM forms WHERE id = $1", id).Scan(&form.Id, &form.User_id, &form.Name, &form.Target_email, &form.Created_at)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to get form", 500)
		return
	}
	response.HttpResponse(w, form, 1, "", 200)
}

func UpdateForm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	var updateForm Form
	_ = json.NewDecoder(r.Body).Decode(&updateForm)

	_, dbErr := database.Db.Exec("UPDATE forms SET name = $1, target_email = $2 WHERE id = $3", updateForm.Name, updateForm.Target_email, id)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to update form", 500)
		return
	}
	response.HttpResponse(w, "", 1, "Form Updated", 200)
}

func DeleteForm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := (params["id"])
	if id == "" {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	_, dbErr := database.Db.Exec("DELETE FROM forms WHERE id = $1", id)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to delete form", 500)
		return
	}
	response.HttpResponse(w, "", 1, "Form Deleted", 200)
}
