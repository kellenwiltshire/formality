package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kellenwiltshire/formality/internal/store"
	"github.com/kellenwiltshire/formality/internal/util"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("missing user in request") // bad actor call
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("formality_auth")
		if err != nil {
			if err != http.ErrNoCookie {
				util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "no token received"})
				return
			}
			util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal service error"})
			return
		}

		token := c.Value
		user, err := um.UserStore.GetUserToken(token)
		if err != nil {
			util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "invalid token"})
			return
		}

		if user == nil {
			util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "token expired or invalid"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "you must be logged in to access this route"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) AuthenticateAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("formality_auth")
		if err != nil {
			if err == http.ErrNoCookie {
				util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "token expired or invalid"})
				return
			}
			util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal service error"})
			return
		}

		token := c.Value
		user, err := um.UserStore.GetAdminToken(token)
		if err != nil {
			fmt.Printf("ERROR: %v \n", err)
			util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "invalid token"})
			return
		}

		if user == nil {
			util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "token expired or invalid"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "you must be logged in to access this route"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
