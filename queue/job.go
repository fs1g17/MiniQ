package queue

import "errors"

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

type Job[T any] struct {
	ID       int
	Status   JobStatus
	Data     T
	Attempts int
}

func (j *Job[T]) UpdateStatus(js JobStatus) {
	j.Status = js
}

func NewJob[T any](data T) *Job[T] {
	return &Job[T]{
		Status:   Queued,
		Data:     data,
		Attempts: 0,
	}
}
