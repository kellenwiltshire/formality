package api

import (
	"database/sql"
	"encoding/json"
	"formality/internal/middleware"
	"formality/internal/store"
	"formality/internal/util"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type registerSubmissionRequest struct {
	Payload string `json:"payload"`
}

type SubmissionHandler struct {
	submissionsStore store.SubmissionsStore
	logger           *log.Logger
}

func NewSubmissionHandler(submissionsStore store.SubmissionsStore, logger *log.Logger) *SubmissionHandler {
	return &SubmissionHandler{
		submissionsStore: submissionsStore,
		logger:           logger,
	}
}

func (h *SubmissionHandler) HandleCreateSubmission(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		h.logger.Printf("Invalid Id param")
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid id"})
		return
	}

	var req registerSubmissionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding create submission request %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid submission payload"})
		return
	}

	submission := &store.Submission{
		Payload: req.Payload,
	}

	err = h.submissionsStore.CreateSubmission(submission)
	if err != nil {
		h.logger.Printf("ERROR: registering submission %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusCreated, util.Envelope{"submission": submission})
}

func (h *SubmissionHandler) HandleGetFormSubmissions(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		h.logger.Printf("Invalid Id param")
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid id"})
		return
	}

	err := h.submissionsStore.DoesUserOwnForm(int64(user.Id), idParam)
	if err != nil {
		h.logger.Printf("Unable to get form ownership %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "Unable to get form ownership"})
		return
	}

	submissions, err := h.submissionsStore.GetFormSubmissions(idParam)
	if err != nil {
		h.logger.Printf("ERROR: getSubmissions: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"submissions": submissions})
}

func (h *SubmissionHandler) HandleGetFormSubmissionById(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	idParam := chi.URLParam(r, "form_id")
	if idParam == "" {
		h.logger.Printf("Invalid Id param")
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid id"})
		return
	}

	err := h.submissionsStore.DoesUserOwnForm(int64(user.Id), idParam)
	if err != nil {
		h.logger.Printf("Unable to get form ownership %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "Unable to get form ownership"})
		return
	}

	submissionParam := chi.URLParam(r, "submission_id")
	if idParam == "" {
		h.logger.Printf("Invalid Id param")
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid id"})
		return
	}

	submissionId, err := strconv.ParseInt(submissionParam, 10, 64)
	if err != nil {
		h.logger.Printf("Unable to parse ID %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "Unable to parse id"})
		return
	}

	submission, err := h.submissionsStore.GetFormSubmissionById(submissionId)
	if err != nil {
		h.logger.Printf("ERROR: getSubmission: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"submission": submission})
}

func (h *SubmissionHandler) HandleDeleteFormSubmission(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	idParam := chi.URLParam(r, "form_id")
	if idParam == "" {
		h.logger.Printf("Invalid Id param")
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid id"})
		return
	}

	err := h.submissionsStore.DoesUserOwnForm(int64(user.Id), idParam)
	if err != nil {
		h.logger.Printf("Unable to get form ownership %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "Unable to get form ownership"})
		return
	}

	submissionParam := chi.URLParam(r, "submission_id")
	if idParam == "" {
		h.logger.Printf("Invalid Id param")
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid id"})
		return
	}

	submissionId, err := strconv.ParseInt(submissionParam, 10, 64)
	if err != nil {
		h.logger.Printf("Unable to parse ID %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "Unable to parse id"})
		return
	}

	err = h.submissionsStore.DeleteSubmission(submissionId)
	if err == sql.ErrNoRows {
		http.Error(w, "submission not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error deleting submission", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
