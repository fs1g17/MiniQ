package queue

type JobStatus int

const (
	Queued JobStatus = iota
	Processing
	Completed
	Failed
)

var statusName = map[JobStatus]string{
	Queued:     "queued",
	Processing: "processing",
	Completed:  "completed",
	Failed:     "failed",
}

func (js JobStatus) String() string {
	return statusName[js]
}

type Job[T any] struct {
	Name   string
	Status JobStatus
	Data   T
}
