package repository

type LogMessage map[string]interface{}

func (r *SQLite) StoreLogs(logs []LogMessage) {
	r.db.NamedExec("INSERT INTO logs (message, containerID, timestamp ) VALUES (:message, :containerID, :timestamp)", logs)
}
