package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const (
	dbPath = "data/sqlite.db"
)

type SQLite struct {
	db *sqlx.DB
}

func New() *SQLite {
	if _, err := os.Stat("data"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("data", os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		log.Fatalln(err)
	}

	migrate(db)

	return &SQLite{db: db}
}

func migrate(db *sqlx.DB) {
	localTableSchema := `
	CREATE TABLE IF NOT EXISTS local(
	    api_key BLOB,
	    containers TEXT
	);
`

	logsSchema := `
	CREATE TABLE IF NOT EXISTS logs(
	    message TEXT,
	    containerID TEXT,
	    timestamp INTEGER
	);
`

	db.MustExec(localTableSchema)
	db.MustExec(logsSchema)
}
