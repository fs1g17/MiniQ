package queue

import "sync"

type Queue[T any] struct {
	Jobs []*Job[T]
	mu   sync.Mutex
}

func (q *Queue[T]) Enqueue(job *Job[T]) {
	q.Jobs = append(q.Jobs, job)
}

func (q *Queue[T]) Dequeue() *Job[T] {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.Jobs) == 0 {
		return nil
	}

	job := q.Jobs[0]
	q.Jobs = q.Jobs[1:]

	return job
}
