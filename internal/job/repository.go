package job

import (
	"async-job-platform/internal/db"
	"encoding/json"
)

func Create(jobType string, payload map[string]interface{}) (int64, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	var id int64
	query := `INSERT INTO jobs (job_type,payload) VALUES ($1,$2) RETURNING id`
	err = db.DB.QueryRow(query, jobType, payloadBytes).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil

}

type Job struct {
	ID         int64                  `json:"id"`
	JobType    string                 `json:"job_type"`
	Status     string                 `json:"status"`
	RetryCount int                    `json:"retry_count"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
}

func GetByID(id int64) (*Job, error) {
	var job Job
	query := `SELECT id,job_type,status,retry_count FROM jobs WHERE id=$1`
	err := db.DB.QueryRow(query, id).Scan(&job.ID, &job.JobType, &job.Status, &job.RetryCount)
	if err != nil {
		return nil, err

	}
	return &job, nil
}

func UpdateStatus(id int64, status string) error {
	query := `UPDATE jobs SET status=$1 WHERE id=$2`
	_, err := db.DB.Exec(query, status, id)
	return err
}

func IncrementRetryCount(id int64) error {
	query := `UPDATE jobs SET retry_count=retry_count+1 WHERE id=$1`
	_, err := db.DB.Exec(query, id)
	return err
}

func GetJob(id int64) (*Job, error) {
	var j Job
	query := `SELECT id, job_type,status,retry_count FROM jobs WHERE id=$1`
	err := db.DB.QueryRow(query, id).Scan(
		&j.ID,
		&j.JobType,
		&j.Status,
		&j.RetryCount,
	)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

type Metrics struct {
	Queued     int `json:"queued"`
	Processing int `json:"processing"`
	Completed  int `json:"completed"`
	Failed     int `json:"failed"`
}

func GetMetrics() (*Metrics, error) {
	var m Metrics
	query := `
		SELECT
			COUNT(*) FILTER (WHERE status = 'queued'),
			COUNT(*) FILTER (WHERE status = 'processing'),
			COUNT(*) FILTER (WHERE status = 'completed'),
			COUNT(*) FILTER (WHERE status = 'failed')
		FROM jobs
	`

	err := db.DB.QueryRow(query).Scan(
		&m.Queued,
		&m.Processing,
		&m.Completed,
		&m.Failed,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
