package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kellenwiltshire/formality/internal/middleware"
	"github.com/kellenwiltshire/formality/internal/store"
	"github.com/kellenwiltshire/formality/internal/util"
)

type registerUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.Role == "" {
		return errors.New("role is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding create user request %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request payload"})
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Email: req.Email,
		Role:  req.Role,
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: hashing password %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR: registering user %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusCreated, util.Envelope{"user": user})
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	util.WriteJSON(w, http.StatusOK, util.Envelope{"user": user})
}

func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	existingUser, err := h.userStore.GetUser(int64(user.Id))
	if err != nil {
		h.logger.Printf("ERROR: GetUser: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	if existingUser == nil {
		http.NotFound(w, r)
		return
	}

	var updateUser struct {
		Email    *string `json:"email"`
		Role     *string `json:"role"` // TODO update this so that non-admin can't change their role...
		Password *string `json:"password"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		h.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request payload"})
		return
	}

	if updateUser.Email != nil {
		existingUser.Email = *updateUser.Email
	}

	if updateUser.Role != nil {
		existingUser.Role = *updateUser.Role
	}

	if updateUser.Password != nil {
		err = existingUser.PasswordHash.Set(*updateUser.Password)
		if err != nil {
			h.logger.Printf("ERROR: hashing password %v", err)
			util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
			return
		}
	}

	err = h.userStore.UpdateUser(existingUser)
	if err != nil {
		h.logger.Printf("ERROR: updatingUser: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"user": existingUser.Id})
}

func (h *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := util.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("Error: ReadParamId: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid user id"})
		return
	}

	err = h.userStore.DeleteUser(userId)
	if err == sql.ErrNoRows {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userStore.GetAllUsers()
	if err != nil {
		h.logger.Printf("Error: getAllUsers %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "Error getting all users"})
		return
	}
	util.WriteJSON(w, http.StatusOK, util.Envelope{"users": users})
}

func (h *UserHandler) HandleCreateAdminUser() error {
	numUsers, err := h.userStore.GetNumberAdminUsers()
	if err != nil {
		return err
	}

	if numUsers == nil {
		return fmt.Errorf("Could not determine number of users")
	}

	if *numUsers == 0 {
		adminPass := os.Getenv("ADMIN_PASS")
		if adminPass == "" {
			return fmt.Errorf("Must provide a default Admin Password")
		}

		user := &store.User{
			Email: "admin@email.com",
			Role:  "admin",
		}

		err = user.PasswordHash.Set(adminPass)
		if err != nil {
			return err
		}

		err = h.userStore.CreateUser(user)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Already an Admin User")
	}
	return nil

}
