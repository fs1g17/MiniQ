package store

import (
	"database/sql"
	"encoding/json"

	"github.com/fs1g17/MiniQ/queue"
)

type JobStore[T any] struct {
	db *sql.DB
}

func NewJobStore[T any](db *sql.DB) *JobStore[T] {
	return &JobStore[T]{db: db}
}

func (js *JobStore[T]) InsertJob(job *queue.Job[T]) error {
	query := `
		INSERT INTO jobs (status, data, attempts)
		VALUES ($1, $2::jsonb, $3)
		RETURNING id;
	`

	m, err := json.Marshal(job.Data)
	if err != nil {
		return err
	}

	var jobId int
	err = js.db.QueryRow(query, job.Status.String(), m, job.Attempts).Scan(&jobId)
	if err != nil {
		return err
	}

	job.ID = jobId
	return nil
}

func (js *JobStore[T]) GetJob(id int) (*queue.Job[T], error) {
	query := `
	SELECT id, status, data, attempts
	FROM jobs
	WHERE id = $1;
	`

	var jobStatus string
	var data string

	var job queue.Job[T]
	err := js.db.QueryRow(query, id).Scan(
		&job.ID,
		&jobStatus,
		&data,
		&job.Attempts,
	)

	actualJobStatus, err := queue.GetJobStatus(jobStatus)
	if err != nil {
		return nil, err
	}
	job.Status = actualJobStatus

	return &job, nil
}
