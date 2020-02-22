package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	psqlUser = "PSQL_USER"
	psqlPass = "PSQL_PASS"
	psqlDB   = "PSQL_DB"

	connFormat = "host=localhost port=5432 user=%v password=%v dbname=%v sslmode=disable"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage() (*Storage, error) {
	connF := fmt.Sprintf(connFormat, os.Getenv(psqlUser), os.Getenv(psqlPass), os.Getenv(psqlDB))
	database, err := sql.Open("postgres", connF)

	if err != nil {
		return nil, err
	}

	err = database.Ping()

	if err != nil {
		return nil, err
	}

	return &Storage{
		DB: database,
	}, nil
}
