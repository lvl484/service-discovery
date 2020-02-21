package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/casbin/casbin"
)

const IdleTimeout = 30 * time.Minute

var sessionManager *scs.SessionManager

func main() {
	addr := flag.String("a", ":8080", "port")
	flag.Parse()

	users := newUsers()
	defer users.db.Close()

	data := NewData()

	authEnforce, err := casbin.NewEnforcer("./auth.conf", "policy.csv")

	if err != nil {
		log.Fatal(err)
	}

	sessionManager = scs.New()
	sessionManager.IdleTimeout = IdleTimeout
	sessionManager.Store = postgresstore.New(users.db)
	//TODO: make connects via https
	//sessionManager.Cookie.Secure = true
	mainRouter := newRouter(data, users)

	if err := http.ListenAndServe(
		*addr, sessionManager.LoadAndSave(Authorizer(authEnforce)(mainRouter)),
	); err != nil {
		log.Fatal(err)
	}
}
