package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	addr := flag.String("a", ":8080", "port")
	flag.Parse()

	data := newData()

	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/", data.Add).Methods(http.MethodPost)
	mainRouter.HandleFunc("/", data.GetIDs).Methods(http.MethodGet)
	mainRouter.HandleFunc("/{ID}", data.Get).Methods(http.MethodGet)
	mainRouter.HandleFunc("/{ID}", data.Update).Methods(http.MethodPut)
	mainRouter.HandleFunc("/{ID}", data.Delete).Methods(http.MethodDelete)

	if err := http.ListenAndServe(*addr, mainRouter); err != nil {
		log.Fatal(err)
	}
}
