package queue

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
