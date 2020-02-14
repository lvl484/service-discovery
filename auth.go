package main

import (
	"net/http"

	"github.com/casbin/casbin"
)

// Authorizer is a middleware for authorization
func Authorizer(e *casbin.Enforcer, users *Users) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			role := sessionManager.GetString(r.Context(), "role")
			if role == "" {
				role = "guest"
			} else if len(role) > 0 {
				uid := sessionManager.GetString(r.Context(), "userID")
				//uid := "1"
				exists := users.Exists(uid)
				if !exists {
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
			} else {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		return http.HandlerFunc(fn)
	}
}
