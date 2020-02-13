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

	mainRouter := mux.NewRouter()

	if err := http.ListenAndServe(*addr, mainRouter); err != nil {
		log.Fatal(err)
	}
}
