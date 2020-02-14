package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
)

const IdleTimeout = 10

var sessionManager *scs.SessionManager

func main() {
	addr := flag.String("a", ":8080", "port")
	flag.Parse()

	users := newUsers()
	defer users.db.Close()

	data := newData()

	authEnforce, err := casbin.NewEnforcer("./auth.conf", "policy.csv")

	if err != nil {
		log.Fatal(err)
	}

	sessionManager = scs.New()
	sessionManager.IdleTimeout = IdleTimeout * time.Minute

	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/list", data.Add).Methods(http.MethodPost)
	mainRouter.HandleFunc("/list", data.GetIDs).Methods(http.MethodGet)
	mainRouter.HandleFunc("/list/{ID}", data.Get).Methods(http.MethodGet)
	mainRouter.HandleFunc("/list/{ID}", data.Update).Methods(http.MethodPut)
	mainRouter.HandleFunc("/list/{ID}", data.Delete).Methods(http.MethodDelete)
	mainRouter.HandleFunc("/login", loginHandler(users))

	if err := http.ListenAndServe(
		*addr, sessionManager.LoadAndSave(Authorizer(authEnforce, users)(mainRouter)),
	); err != nil {
		log.Fatal(err)
	}
}
