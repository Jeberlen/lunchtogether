package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDB() {

	connStr := "postgres://postgres:postgres@localhost/lunchtogetherdev?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
 		log.Panic(err)
	}
	Db = db
}

func CloseDB() error {
	return Db.Close()
}

