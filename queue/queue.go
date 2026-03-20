package queue

import "sync"

type Queue[T any] struct {
	Jobs    []*Job[T]
	mu      sync.Mutex
	Workers []*Worker[T]
}

func (q *Queue[T]) tryQueue() {
	numberJobs := len(q.Jobs)
	q.mu.Unlock()
	// no jobs to process
	if (numberJobs) == 0 {
		return
	}
	// there are jobs, look for available worker
	var availableWorker *Worker[T] = nil
	for _, worker := range q.Workers {
		if worker.Status == Busy {
			continue
		}
		availableWorker = worker
		break
	}
	if availableWorker != nil {
		go availableWorker.Perform()
	}
}

func (q *Queue[T]) Try() {
	q.mu.Lock()
	defer q.tryQueue()
}

func (q *Queue[T]) Enqueue(job *Job[T]) {
	q.mu.Lock()
	defer q.tryQueue()
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
