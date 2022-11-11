package repository

func (r *SQLite) StoreAPIKey(key []byte) {
	r.db.Exec("INSERT INTO local (api_key) VALUES ($1)", key)
}

func (r *SQLite) GetAPIKey() []byte {
	var key []byte
	r.db.Get(&key, "SELECT api_key FROM local LIMIT 1;")

	return key
}
