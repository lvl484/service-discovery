package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/casbin/casbin"

	"github.com/lvl484/service-discovery/servicetrace"
	"github.com/lvl484/service-discovery/storage"
)

const SessionIdleTimeout = 30 * time.Minute

var sessionManager *scs.SessionManager

func main() {
	addr := flag.String("a", ":8080", "port")
	flag.Parse()

	storage, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer storage.DB.Close()

	userStorage := newUserStorage(storage.DB)
	data := NewData()
	services := servicetrace.NewServices()

	authEnforce, err := casbin.NewEnforcer("./auth.conf", "policy.csv")

	if err != nil {
		log.Fatal(err)
	}

	sessionManager = scs.New()
	sessionManager.IdleTimeout = SessionIdleTimeout
	sessionManager.Store = postgresstore.New(storage.DB)
	//TODO: make connects via https
	//sessionManager.Cookie.Secure = true
	mainRouter := newRouter(data, userStorage, services)

	if err := http.ListenAndServe(
		*addr, sessionManager.LoadAndSave(Authorizer(authEnforce)(mainRouter)),
	); err != nil {
		log.Fatal(err)
	}
}
