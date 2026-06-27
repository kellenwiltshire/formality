package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"formality/internal/middleware"
	"formality/internal/service"
	"formality/internal/store"
	"formality/internal/util"
	"log"
	"net/http"
	"regexp"
)

type registerSmtpRequest struct {
	UserID         int    `json:"user_id"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	EncryptionType string `json:"encryption_type"`
	RecipientEmail string `json:"recipient_email"`
	SenderEmail    string `json:"sender_email"`
}

type SmtpHandler struct {
	smtpStore       store.SmtpStore
	sendMailService service.SendMailService
	logger          *log.Logger
}

func NewSmtpHandler(smtpStore store.SmtpStore, sendMailService service.SendMailService, logger *log.Logger) *SmtpHandler {
	return &SmtpHandler{
		smtpStore:       smtpStore,
		sendMailService: sendMailService,
		logger:          logger,
	}
}

func (h *SmtpHandler) validateRegisterRequest(req *registerSmtpRequest) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if req.Username == "" {
		return errors.New("username is required")
	}

	if req.RecipientEmail == "" {
		return errors.New("recipient email is required")
	}

	if !emailRegex.MatchString(req.RecipientEmail) {
		return errors.New("invalid recipient email format")
	}

	if req.SenderEmail == "" {
		return errors.New("sender email is required")
	}

	if !emailRegex.MatchString(req.SenderEmail) {
		return errors.New("invalid sender email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	if req.Host == "" {
		return errors.New("host is required")
	}

	if req.Port == 0 {
		return errors.New("Port is required")
	}

	return nil
}

func (h *SmtpHandler) HandleCreateSmtpSettings(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	existingSmtp, err := h.smtpStore.GetSmtpSettings(int64(user.Id))
	if err != nil {
		h.logger.Printf("ERROR: GetSmtp: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	if existingSmtp != nil {
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "smtp settings already exist"})
		return
	}

	var req registerSmtpRequest

	err = json.NewDecoder(r.Body).Decode(&req)
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

	smtp := &store.Smtp{
		UserID:         user.Id,
		Host:           req.Host,
		Port:           req.Port,
		Username:       req.Username,
		EncryptionType: req.EncryptionType,
		RecipientEmail: req.RecipientEmail,
		SenderEmail:    req.SenderEmail,
	}

	err = smtp.PasswordEncrypted.EncryptAES(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: hashing password %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	err = h.smtpStore.CreateSmtpSettings(smtp)
	if err != nil {
		h.logger.Printf("ERROR: registering user %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusCreated, util.Envelope{"smtp": smtp})
}

func (h *SmtpHandler) HandleGetSMTPSettings(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	smtp, err := h.smtpStore.GetSmtpSettings(int64(user.Id))
	if err != nil {
		h.logger.Printf("ERROR: GetSMTP: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"smtp": smtp})
}

func (h *SmtpHandler) HandleUpdateSmtpSettings(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	existingSmtp, err := h.smtpStore.GetSmtpSettings(int64(user.Id))
	if err != nil {
		h.logger.Printf("ERROR: GetSmtp: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	if existingSmtp == nil {
		http.NotFound(w, r)
		return
	}

	var updateSmtp struct {
		Host           *string `json:"host"`
		Port           *int    `json:"port"`
		Username       *string `json:"username"`
		Password       *string `json:"password"`
		EncryptionType *string `json:"encryption_type"`
		RecipientEmail *string `json:"recipient_email"`
		SenderEmail    *string `json:"sender_email"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateSmtp)
	if err != nil {
		h.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request payload"})
		return
	}

	if updateSmtp.Host != nil {
		existingSmtp.Host = *updateSmtp.Host
	}

	if updateSmtp.Port != nil {
		existingSmtp.Port = *updateSmtp.Port
	}

	if updateSmtp.Username != nil {
		existingSmtp.Username = *updateSmtp.Username
	}

	if updateSmtp.EncryptionType != nil {
		existingSmtp.EncryptionType = *updateSmtp.EncryptionType
	}

	if updateSmtp.RecipientEmail != nil {
		existingSmtp.RecipientEmail = *updateSmtp.RecipientEmail
	}

	if updateSmtp.SenderEmail != nil {
		existingSmtp.SenderEmail = *updateSmtp.SenderEmail
	}

	if updateSmtp.Password != nil {
		err = existingSmtp.PasswordEncrypted.EncryptAES(*updateSmtp.Password)
		if err != nil {
			h.logger.Printf("ERROR: hashing password %v", err)
			util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
			return
		}
	}

	err = h.smtpStore.UpdateSmtpSettings(existingSmtp)
	if err != nil {
		h.logger.Printf("ERROR: updatingSmtp: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"smtp": existingSmtp.Id})
}

func (h *SmtpHandler) HandleDeleteSmtpSetting(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	err := h.smtpStore.DeleteSmtpSettings(int64(user.Id))
	if err == sql.ErrNoRows {
		http.Error(w, "smtp not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error deleting smtp", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SmtpHandler) HandleTestEmail(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	payload := "This is a test email for Formality!"

	err := h.sendMailService.TestSendMail(int64(user.Id), payload)
	if err != nil {
		h.logger.Printf("ERROR: TestEmail: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"test": "success"})
}
