package queue

import (
	"fmt"
	"sync"
)

type WorkerStatus int

const (
	Busy WorkerStatus = iota
	Idle
)

var workerStatusName = map[WorkerStatus]string{
	Busy: "busy",
	Idle: "idle",
}

func (ws WorkerStatus) String() string {
	return workerStatusName[ws]
}

type Worker[T any] struct {
	ID         int
	Work       func(T) error
	LogChannel chan string
	JobChannel chan string
	Status     WorkerStatus
	mu         sync.Mutex
}

func (w *Worker[T]) SetStatus(ws WorkerStatus) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Status = ws
}

func (w *Worker[T]) GetStatus() WorkerStatus {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Status
}

func (w *Worker[T]) Perform(job *Job[T]) {
	w.SetStatus(Busy)
	defer func() {
		w.SetStatus(Idle)
		w.JobChannel <- fmt.Sprintf("WORKER_FREED: %d", w.ID)
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered worker %d jobName %s reason %v\n", w.ID, job.Name, r)
			job.UpdateStatus(Failed)
		}
	}()

	job.UpdateStatus(Processing)
	w.LogChannel <- fmt.Sprintf("job %s initiated by worker %d", job.Name, w.ID)
	err := w.Work(job.Data)
	if err != nil {
		job.UpdateStatus(Failed)
	} else {
		job.UpdateStatus(Completed)
	}
}
