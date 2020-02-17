package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MysqlUser = "MYSQL_USER"
	MysqlPass = "MYSQL_PASS"
	MysqlDb   = "MYSQL_DB"
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

func (u *Users) Exists(id string) (bool, error) {
	var user User

	userRow := u.db.QueryRow("SELECT ID FROM User WHERE ID=?", id)

	err := userRow.Scan(&user.ID)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (u *Users) FindByCredentials(name, pass string) (*User, error) {
	var user User

	userRow := u.db.QueryRow("SELECT ID, ROLE FROM User WHERE Username=? and Password=?", name, pass)
	err := userRow.Scan(&user.ID, &user.Role)

	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func newUsers() *Users {
	s := fmt.Sprintf("%v:%v@/%v", os.Getenv(MysqlUser), os.Getenv(MysqlPass), os.Getenv(MysqlDb))
	database, err := sql.Open("mysql", s)

	if err != nil {
		log.Println(err)
	}

	return &Users{
		db: database,
	}
}
