package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/lvl484/service-discovery/encodepass"
)

const (
	psqlUser = "PSQL_USER"
	psqlPass = "PSQL_PASS"
	psqlDB   = "PSQL_DB"

	connFormat = "host=localhost port=5432 user=%v password=%v dbname=%v sslmode=disable"
)

type User struct {
	ID       string
	Username string
	Password string
	Role     string
}

type Users struct {
	conf *encodepass.PasswordConfig
	db   *sql.DB
}

func (u *Users) Register(user *User) error {
	stmt := `
		INSERT INTO users(ID, Username, Password, Role) VALUES($1,$2,$3,$4) `
	_, err := u.db.Exec(
		stmt,
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *Users) FindByCredentials(name, pass string) (*User, error) {
	var user User

	userRow := u.db.QueryRow("SELECT * FROM users WHERE username=$1", name)
	err := userRow.Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	res, err := encodepass.ComparePassword(pass, user.Password)

	if err != nil {
		return nil, err
	}

	if !res {
		return nil, nil
	}

	return &user, nil
}

func newUsers() *Users {
	connF := fmt.Sprintf(connFormat, os.Getenv(psqlUser), os.Getenv(psqlPass), os.Getenv(psqlDB))
	database, err := sql.Open("postgres", connF)

	if err != nil {
		log.Fatal(err)
	}

	err = database.Ping()

	if err != nil {
		log.Println(err)
	}

	conf := encodepass.NewPasswordConfig()

	return &Users{
		db:   database,
		conf: conf,
	}
}
