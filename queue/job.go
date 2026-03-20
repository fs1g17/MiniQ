package queue

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
	Name   string
	Status JobStatus
	Data   T
}

func (j *Job[T]) UpdateStatus(js JobStatus) {
	j.Status = js
}
