package queue

import "sync"

type Queue[T any] struct {
	jobs []*Job[T]
	mu   sync.Mutex
}

func (q *Queue[T]) enqueue(job *Job[T]) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, job)
}

func (q *Queue[T]) dequeue() *Job[T] {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.jobs) == 0 {
		return nil
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]

	return job
}
