package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/kellenwiltshire/formality/internal/middleware"
	"github.com/kellenwiltshire/formality/internal/store"
	"github.com/kellenwiltshire/formality/internal/util"

	"github.com/go-chi/chi/v5"
)

type createFormRequest struct {
	UserId      int    `json:"user_id"`
	Name        string `json:"name"`
	TargetEmail string `json:"target_email"`
}

type FormHandler struct {
	formStore store.FormStore
	logger    *log.Logger
}

func NewFormHandler(formStore store.FormStore, logger *log.Logger) *FormHandler {
	return &FormHandler{
		formStore: formStore,
		logger:    logger,
	}
}

func (h *FormHandler) validateRegisterRequest(req *createFormRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}

	if req.TargetEmail == "" {
		return errors.New("email is required")
	}

	return nil
}

func (h *FormHandler) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
	var req createFormRequest
	user := middleware.GetUser(r)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding create form request %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request payload"})
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": err.Error()})
		return
	}

	form := &store.Form{
		UserId:      user.Id,
		Name:        req.Name,
		TargetEmail: req.TargetEmail,
	}

	err = h.formStore.CreateForm(form)
	if err != nil {
		h.logger.Printf("ERROR: registering form %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusCreated, util.Envelope{"form": form})
}

func (h *FormHandler) HandleGetForm(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	formId := chi.URLParam(r, "form_id")
	if formId == "" {
		h.logger.Panicf("Error: decodeParam %v", formId)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid id"})
		return
	}

	form, err := h.formStore.GetForm(formId, int64(user.Id))
	if err != nil {
		h.logger.Printf("ERROR: GetForm: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"form": form})
}

func (h *FormHandler) HandleUpdateForm(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	formId := chi.URLParam(r, "form_id")
	if formId == "" {
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid id"})
	}

	existingForm, err := h.formStore.GetForm(formId, int64(user.Id))
	if err != nil {
		h.logger.Printf("ERROR: GetForm: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	if existingForm == nil {
		http.NotFound(w, r)
		return
	}

	var updateForm struct {
		Name        *string `json:"name"`
		TargetEmail *string `json:"target_email"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateForm)
	if err != nil {
		h.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request payload"})
		return
	}

	if updateForm.Name != nil {
		existingForm.Name = *updateForm.Name
	}

	if updateForm.TargetEmail != nil {
		existingForm.TargetEmail = *updateForm.TargetEmail
	}

	err = h.formStore.UpdateForm(existingForm)
	if err != nil {
		h.logger.Printf("ERROR: updatingForm: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"form": existingForm.Id})
}

func (h *FormHandler) HandleDeleteForm(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	formId := chi.URLParam(r, "form_id")
	if formId == "" {
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid id"})
	}

	err := h.formStore.DeleteForm(formId, int64(user.Id))
	if err == sql.ErrNoRows {
		http.Error(w, "form not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error deleting form", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FormHandler) HandleGetAllFormsForUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	forms, err := h.formStore.GetAllFormsForUser(int64(user.Id))
	if err != nil {
		h.logger.Printf("ERROR: GetAllForms: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "Error getting all forms for user"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"forms": forms})
}
