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

	return &SQLite{db: db}
}

func (r *SQLite) MigrateAgentDatabase() {
	localTableSchema := `
	CREATE TABLE IF NOT EXISTS local(
	    api_key BLOB,
	    containers TEXT DEFAULT ''
	);
`

	logsSchema := `
	CREATE TABLE IF NOT EXISTS logs(
	    message TEXT,
	    containerID TEXT,
	    timestamp INTEGER
	);
`

	r.db.MustExec(localTableSchema)
	r.db.MustExec(logsSchema)
}

func (r *SQLite) MigrateClientDatabase() {
	nodesSchema := `
	CREATE TABLE IF NOT EXISTS nodes(
		location TEXT NOT NULL,
		api_key TEXT NOT NULL,
		name TEXT NOT NULL UNIQUE
	);
`

	r.db.MustExec(nodesSchema)
}
