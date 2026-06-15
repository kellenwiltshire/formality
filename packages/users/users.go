package users

import (
	"encoding/json"
	"fmt"
	"formality/packages/database"
	"formality/packages/response"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type NewUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser NewUser
	_ = json.NewDecoder(r.Body).Decode(&newUser)

	if newUser.Email == "" || newUser.Password == "" {
		response.HttpResponse(w, "", 0, "Please provide a username and password", 400)
		return
	}

	hash, _ := HashPassword(newUser.Password)

	_, err := database.Db.Exec("INSERT INTO users (email, password_hash, role) values ($1, $2, $3)", newUser.Email, hash, newUser.Role)
	if err != nil {
		fmt.Println(err)
		response.HttpResponse(w, "", 0, "Unable to create new user", 500)
	} else {
		response.HttpResponse(w, "", 1, "User created", 200)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	var user User
	dbErr := database.Db.QueryRow("SELECT id, email, role FROM users WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.Role)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to get user", 500)
	}
	response.HttpResponse(w, user, 1, "", 200)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	var updateUser NewUser
	_ = json.NewDecoder(r.Body).Decode(&updateUser)

	_, dbErr := database.Db.Exec("UPDATE users SET email = $1, role = $2 WHERE id = $3", updateUser.Email, updateUser.Role, id)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to update user", 500)
	}
	response.HttpResponse(w, "", 1, "User updated", 200)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		response.HttpResponse(w, "", 0, "Invalid ID", 400)
		return
	}

	_, dbErr := database.Db.Exec("DELETE FROM users WHERE id = $1", id)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to delete user", 500)
	}
	response.HttpResponse(w, "", 1, "User Deleted", 200)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	rows, dbErr := database.Db.Query("SELECT id, email, role FROM users")
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to get all users", 500)
	}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Role); err != nil {
			fmt.Println(err)
			response.HttpResponse(w, "", 0, "Unable to get all user rows", 500)
		}
		users = append(users, user)
	}
	response.HttpResponse(w, users, 1, "", 200)
}
