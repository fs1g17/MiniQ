package queue

import (
	"fmt"
	"sync"

	"github.com/fs1g17/MiniQ/store"
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

type Worker struct {
	ID         int
	Work       func(store.AnyData) error
	LogChannel chan string
	JobChannel chan string
	Status     WorkerStatus
	mu         sync.Mutex
}

func (w *Worker) SetStatus(ws WorkerStatus) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Status = ws
}

func (w *Worker) GetStatus() WorkerStatus {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Status
}

func (w *Worker) Perform(job *store.Job) {
	w.SetStatus(Busy)
	defer func() {
		w.SetStatus(Idle)
		w.JobChannel <- fmt.Sprintf("WORKER_FREED: %d", w.ID)
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered worker %d jobID %d reason %v\n", w.ID, job.ID, r)
			job.UpdateStatus(store.Failed)
		}
	}()

	job.UpdateStatus(store.Processing)
	w.LogChannel <- fmt.Sprintf("job %d initiated by worker %d", job.ID, w.ID)
	err := w.Work(job.Data)
	if err != nil {
		job.UpdateStatus(store.Failed)
	} else {
		job.UpdateStatus(store.Completed)
	}
}
