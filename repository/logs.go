package repository

type LogMessage map[string]interface{}

func (r *Repository) StoreLogs(logs []LogMessage) {
	r.db.NamedExec("INSERT INTO logs (message) VALUES (:message)", logs)
}
