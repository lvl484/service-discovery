package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       string
	Username string
	Password string
	Role     string
}

type Users struct {
	db *sql.DB
}

func (u *Users) Exists(id string) bool {
	var user User

	userRow := u.db.QueryRow("SELECT * FROM User WHERE ID=?", id)

	err := userRow.Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return false
	}

	if user.ID == id {
		log.Println("exist")
		return true
	}

	return false
}

func (u *Users) FindByCredentials(name, pass string) (User, error) {
	var user User

	userRow := u.db.QueryRow("SELECT * FROM User WHERE Username=? and Password=?", name, pass)
	err := userRow.Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		return User{}, errors.New("USER_NOT_FOUND")
	}

	return user, nil
}

func newUsers() *Users {
	s := fmt.Sprintf("%v:%v@/%v", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASS"), os.Getenv("MYSQL_DB"))
	database, err := sql.Open("mysql", s)

	if err != nil {
		log.Println(err)
	}

	return &Users{
		db: database,
	}
}
