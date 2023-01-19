package repository

import (
	"strings"
)

type LocalConfig struct {
	ID         int
	APIKey     string `db:"api_key"`
	Containers string
}

func (l *LocalConfig) GetContainers() []string {
	if len(l.Containers) == 0 {
		return []string{}
	}

	return strings.Split(l.Containers, ",")
}

func (r *SQLite) GetConfig() (LocalConfig, error) {
	var config LocalConfig
	if err := r.db.Get(&config, "SELECT * FROM local LIMIT 1"); err != nil {
		return LocalConfig{}, err
	}

	return config, nil
}

func (r *SQLite) StoreAPIKey(key []byte) {
	r.db.Exec("INSERT INTO local (api_key) VALUES ($1)", key)
}

func (r *SQLite) GetAPIKey() []byte {
	var key []byte
	r.db.Get(&key, "SELECT api_key FROM local LIMIT 1;")

	return key
}

func (r *SQLite) StoreContainers(containers string) {
	r.db.Exec("UPDATE local SET containers = $1 WHERE id = 1", containers)
}
