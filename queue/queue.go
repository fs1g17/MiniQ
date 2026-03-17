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

type Job struct {
	Name   string
	Status JobStatus
}

type Queue struct {
	Jobs []*Job
}

func (q *Queue) Enqueue(job *Job) {
	q.Jobs = append(q.Jobs, job)
}

func (q *Queue) Dequeue() *Job {
	job := q.Jobs[0]
	q.Jobs = q.Jobs[1:]

	return job
}
