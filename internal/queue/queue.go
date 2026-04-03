package queue

import (
	"sync"

	"github.com/fs1g17/MiniQ/internal/store"
)

type Queue struct {
	jobs []*store.Job
	mu   sync.Mutex
}

func (q *Queue) enqueue(job *store.Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, job)
}

func (q *Queue) dequeue() *store.Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.jobs) == 0 {
		return nil
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]

	return job
}
