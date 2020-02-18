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

	ErrWrongCredentials = "WRONG_CREDENTIALS"
	ErrRenewToken       = "RENEW_TOKEN_ERR"
	ErrUserExists       = "User_Exists"
	Success             = "SUCCESS"
)

// Authorizer is a middleware for authorization
func Authorizer(e *casbin.Enforcer, users *Users) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			role := sessionManager.GetString(r.Context(), UserRole)
			if role == "" {
				role = DefaultRole
			} else if len(role) > 0 {
				uid := sessionManager.GetString(r.Context(), UserID)
				pass := sessionManager.GetString(r.Context(), Password)
				exists, err := users.Exists(uid, pass)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				} else if !exists {
					w.WriteHeader(http.StatusForbidden)
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

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if user == nil {
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

		if len(name) == 0 && len(pass) == 0 {
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

		_, err = users.db.Exec(
			"INSERT INTO User(ID, Username, Password, Role) VALUES(?,?,?,?)",
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Role,
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}
