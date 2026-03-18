package queue

type Queue[T any] struct {
	Jobs []*Job[T]
}

func (q *Queue[T]) Enqueue(job *Job[T]) {
	q.Jobs = append(q.Jobs, job)
}

func (q *Queue[T]) Dequeue() *Job[T] {
	job := q.Jobs[0]
	q.Jobs = q.Jobs[1:]

	return job
}
