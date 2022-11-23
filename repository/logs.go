package repository

type LogMessage map[string]interface{}

type Log struct {
	Message     string
	ContainerID string `db:"containerID"`
	Timestamp   int64
}

func (r *SQLite) StoreLogs(logs []LogMessage) {
	r.db.NamedExec("INSERT INTO logs (message, containerID, timestamp ) VALUES (:message, :containerID, :timestamp)", logs)
}

func (r *SQLite) GetLastNLogs(n int32) ([]Log, error) {
	logsReversed := []Log{}

	err := r.db.Select(&logsReversed, "SELECT * FROM logs ORDER BY timestamp DESC LIMIT $1", n)
	if err != nil {
		return nil, err
	}

	logs := []Log{}

	for i := len(logsReversed) - 1; i >= 0; i-- {
		logs = append(logs, logsReversed[i])
	}

	return logs, nil
}

func (r *SQLite) GetFirstNLogs(n int32) ([]Log, error) {
	logs := []Log{}

	err := r.db.Select(&logs, "SELECT * FROM logs LIMIT $1", n)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
