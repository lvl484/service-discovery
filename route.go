package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter(data *Data, users *Users) *mux.Router {
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/list", data.Add).Methods(http.MethodPost)
	mainRouter.HandleFunc("/list", data.GetIDs).Methods(http.MethodGet)
	mainRouter.HandleFunc("/list/{ID}", data.Get).Methods(http.MethodGet)
	mainRouter.HandleFunc("/list/{ID}", data.Update).Methods(http.MethodPut)
	mainRouter.HandleFunc("/list/{ID}", data.Delete).Methods(http.MethodDelete)
	mainRouter.HandleFunc("/login", loginHandler(users))
	//	mainRouter.HandleFunc("/join", joinHandler(users))

	return mainRouter
}
