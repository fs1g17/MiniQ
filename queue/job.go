package queue

import "github.com/google/uuid"

type JobStatus int

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

type Job[T any] struct {
	Name     string
	Status   JobStatus
	Data     T
	Attempts int
}

func (j *Job[T]) UpdateStatus(js JobStatus) {
	j.Status = js
}

func NewJob[T any](data T, name *string) *Job[T] {
	var realName string
	if name != nil {
		realName = *name
	} else {
		realName = uuid.NewString()
	}
	return &Job[T]{
		Name:     realName,
		Status:   Queued,
		Data:     data,
		Attempts: 0,
	}
}
