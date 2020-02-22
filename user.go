package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/lvl484/service-discovery/encodepass"
)

type User struct {
	ID       string
	Username string
	Password string
	Role     string
}

type UserStorage struct {
	conf *encodepass.PasswordConfig
	db   *sql.DB
}

func (u *UserStorage) Register(user *User) error {
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

func (u *UserStorage) FindByCredentials(name, pass string) (*User, error) {
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

func (u *UserStorage) FindByUsername(name string) (*User, error) {
	var user User

	userRow := u.db.QueryRow("SELECT username FROM users WHERE username=$1", name)
	err := userRow.Scan(&user.Username)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func newUserStorage(db *sql.DB) *UserStorage {
	conf := encodepass.NewPasswordConfig()

	return &UserStorage{
		db:   db,
		conf: conf,
	}
}
