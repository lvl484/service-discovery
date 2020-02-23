package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lvl484/service-discovery/servicetrace"
)

func newRouter(data *Data, users *UserStorage, services *servicetrace.Services) *mux.Router {
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/list", data.Add).Methods(http.MethodPost)
	mainRouter.HandleFunc("/list", data.GetIDs).Methods(http.MethodGet)
	mainRouter.HandleFunc("/list/{ID}", data.Get).Methods(http.MethodGet)
	mainRouter.HandleFunc("/list/{ID}", data.Update).Methods(http.MethodPut)
	mainRouter.HandleFunc("/list/{ID}", data.Delete).Methods(http.MethodDelete)
	mainRouter.HandleFunc("/list/{ID}/{SERVICE_NAME}", data.GetForService(services)).Methods(http.MethodGet)
	mainRouter.HandleFunc("/service-list", GetListOfServices(*services)).Methods(http.MethodGet)
	mainRouter.HandleFunc("/login", loginHandler(users))
	mainRouter.HandleFunc("/logout", logoutHandler())
	mainRouter.HandleFunc("/join", joinHandler(users))

	return mainRouter
}
