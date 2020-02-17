package main

import (
	"database/sql"
	"net/http"

	"github.com/casbin/casbin"
)

const (
	UserID   = "id"
	UserRole = "role"
	Username = "username"
	Password = "password"

	ErrWrongCredentials = "WRONG_CREDENTIALS"
	ErrRenewToken       = "RENEW_TOKEN_ERR"
	Success             = "SUCCESS"
)

// Authorizer is a middleware for authorization
func Authorizer(e *casbin.Enforcer, users *Users) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			role := sessionManager.GetString(r.Context(), UserRole)
			if role == "" {
				role = "guest"
			} else if len(role) > 0 {
				uid := sessionManager.GetString(r.Context(), UserID)
				exists, err := users.Exists(uid)
				if !exists && err == sql.ErrNoRows {
					w.WriteHeader(http.StatusForbidden)
					return
				} else if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			// casbin enforce
			res, err := e.Enforce(role, r.URL.Path, r.Method)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if res {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusForbidden)
		}

		return http.HandlerFunc(fn)
	}
}

func loginHandler(users *Users) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue(Username)
		pass := r.PostFormValue(Password)
		user, err := users.FindByCredentials(name, pass)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(ErrWrongCredentials))
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// setup session
		if err := sessionManager.RenewToken(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(ErrRenewToken))
			return
		}
		sessionManager.Put(r.Context(), UserID, &user.ID)
		sessionManager.Put(r.Context(), UserRole, &user.Role)
		w.Write([]byte(Success))
	})
}
