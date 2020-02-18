package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lvl484/service-discovery/encodepass"
)

const (
	MysqlUser = "MYSQL_USER"
	MysqlPass = "MYSQL_PASS"
	MysqlDB   = "MYSQL_DB"
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

func (u *Users) Exists(id, pass string) (bool, error) {
	var user User

	passRow := u.db.QueryRow("SELECT Password FROM User WHERE ID=?", id)
	err := passRow.Scan(&user.Password)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	res := encodepass.CompareEncodedPassword(pass, user.Password)

	if !res {
		return false, nil
	}

	return true, nil
}

func (u *Users) FindByCredentials(name, pass string) (*User, error) {
	var user User

	userRow := u.db.QueryRow("SELECT * FROM User WHERE Username=?", name)
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
	s := fmt.Sprintf("%v:%v@/%v", os.Getenv(MysqlUser), os.Getenv(MysqlPass), os.Getenv(MysqlDB))
	database, err := sql.Open("mysql", s)

	if err != nil {
		log.Println(err)
	}

	conf := encodepass.NewPasswordConfig()

	return &Users{
		db:   database,
		conf: conf,
	}
}
