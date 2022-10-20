package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	dbFileName = "db"
)

type Repository struct {
	db *sqlx.DB
}

func New() *Repository {
	db, err := sqlx.Connect("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{db: db}
}
