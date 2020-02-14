package main

import (
	"database/sql"
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
	//userRow := u.db.QueryRow("SELECT * FROM User LIMIT 1")

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

//func FindByCredentials(name string) (User, error) {
//
//	return User{}, errors.New("USER_NOT_FOUND")
//}

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
