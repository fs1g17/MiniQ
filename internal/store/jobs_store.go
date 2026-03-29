package store

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type AnyData map[string]any

func (a AnyData) Value() (driver.Value, error) {
	j, err := json.Marshal(a)
	return j, err
}

func (a *AnyData) Scan(src any) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i any
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*a, ok = i.(map[string]any)
	if !ok {
		return errors.New("Type assertion .(map[string]any) failed.")
	}

	return nil
}

type JobStatus int

var errInvalidJobStatus = errors.New("error string is invalid")

const (
	Queued JobStatus = iota
	Processing
	Completed
	Failed
)

var jobStatusName = map[JobStatus]string{
	Queued:     "queued",
	Processing: "processing",
	Completed:  "completed",
	Failed:     "failed",
}

func (js JobStatus) String() string {
	return jobStatusName[js]
}

func GetJobStatus(status string) (JobStatus, error) {
	switch status {
	case "queued":
		return Queued, nil
	case "processing":
		return Processing, nil
	case "completed":
		return Completed, nil
	case "failed":
		return Failed, nil
	default:
		return Queued, errInvalidJobStatus
	}
}

type Job struct {
	ID        int
	Status    JobStatus
	Data      AnyData
	Attempts  int
	CreatedAt time.Time
}

func (j *Job) UpdateStatus(js JobStatus) {
	j.Status = js
}

type JobStore struct {
	db *sql.DB
}

func NewJobStore(db *sql.DB) *JobStore {
	return &JobStore{db: db}
}

func (js *JobStore) InsertJob(job *Job) error {
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

func (js *JobStore) UpdateJobStatus(jobId int, jobStatus JobStatus) error {
	query := `
	UPDATE jobs
	SET jobStatus = $1 
	WHERE id = $2;
	`

	_, err := js.db.Exec(query, jobStatus.String(), jobId)
	if err != nil {
		return err
	}

	return nil
}

func (js *JobStore) GetJob(id int) (*Job, error) {
	query := `
	SELECT id, status, data, attempts
	FROM jobs
	WHERE id = $1;
	`

	var jobStatus string

	var job Job
	err := js.db.QueryRow(query, id).Scan(
		&job.ID,
		&jobStatus,
		&job.Data,
		&job.Attempts,
	)

	actualJobStatus, err := GetJobStatus(jobStatus)
	if err != nil {
		return nil, err
	}
	job.Status = actualJobStatus

	return &job, nil
}
