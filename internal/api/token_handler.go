package api

import (
	"encoding/json"
	"formality/internal/middleware"
	"formality/internal/store"
	"formality/internal/tokens"
	"formality/internal/util"
	"log"
	"net/http"
	"time"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: createTokenRequest: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request payload"})
		return
	}

	// lets get the user
	user, err := h.userStore.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: GetUserByEmail: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("Error: PasswordHash.Matches %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	if !passwordsDoMatch {
		util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "invalid credentials"})
		return
	}

	var scope string

	if user.Role == tokens.ScopeAdmin {
		scope = tokens.ScopeAdmin
	} else {
		scope = tokens.ScopeAuth
	}

	token, err := h.tokenStore.CreateNewToken(user.Id, 24*time.Hour, scope)
	if err != nil {
		h.logger.Printf("Error: Creating Token %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return

	}

	http.SetCookie(w, &http.Cookie{
		Name:     "formality_auth",
		Value:    token.Plaintext,
		Expires:  token.Expiry,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	util.WriteJSON(w, http.StatusCreated, util.Envelope{"auth_token": token})

}

func (h *TokenHandler) HandleDeleteTokens(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	err := h.tokenStore.DeleteAllTokensForUser(user.Id)
	if err != nil {
		h.logger.Printf("Error: Deleting Tokens %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
