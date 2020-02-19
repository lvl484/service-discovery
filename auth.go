package main

import (
	"net/http"

	"github.com/casbin/casbin"
	"github.com/google/uuid"
	"github.com/lvl484/service-discovery/encodepass"
)

const (
	UserID              = "id"
	UserRole            = "role"
	Username            = "username"
	Password            = "password"
	DefaultRole         = "guest"
	DefaultRegisterRole = "user"
	EmptyRole           = ""

	ErrWrongCredentials = "WRONG_CREDENTIALS"
	ErrRenewToken       = "RENEW_TOKEN_ERR"
	ErrUserExists       = "User_Exists"
	Success             = "SUCCESS"
)

// Authorizer is a middleware for authorization
func Authorizer(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			role := sessionManager.GetString(r.Context(), UserRole)
			if role == EmptyRole {
				role = DefaultRole
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

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if user == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(ErrWrongCredentials))
			return
		}

		// setup session
		if err := sessionManager.RenewToken(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(ErrRenewToken))
			return
		}
		sessionManager.Put(r.Context(), UserID, user.ID)
		sessionManager.Put(r.Context(), UserRole, user.Role)
		sessionManager.Put(r.Context(), Username, user.Username)
		sessionManager.Put(r.Context(), Password, user.Password)
		w.Write([]byte(Success))
	})
}

func logoutHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := sessionManager.Destroy(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func joinHandler(users *Users) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user User

		name := r.PostFormValue(Username)
		pass := r.PostFormValue(Password)
		res, err := users.FindByCredentials(name, pass)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if res != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(ErrUserExists))
			return
		}

		if len(name) == 0 || len(pass) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		encodedPass, err := encodepass.EncodePassword(users.conf, pass)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.ID = uuid.New().String()
		user.Username = name
		user.Password = encodedPass
		user.Role = DefaultRegisterRole

		err = users.Register(&user)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}
