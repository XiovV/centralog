package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	dbFileName = "db.db"
)

type SQLite struct {
	db *sqlx.DB
}

func New() *SQLite {
	db, err := sqlx.Connect("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
	}

	return &SQLite{db: db}
}
