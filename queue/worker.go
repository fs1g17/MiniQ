package queue

import (
	"fmt"
	"time"
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
	ID      int
	Work    func(T) error
	Queue   *Queue[T]
	Channel chan string
	Status  WorkerStatus
}

func (w *Worker[T]) Perform() {
	time.Sleep(200 * time.Millisecond)
	job := w.Queue.Dequeue()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered worker %d jobName %s reason %v\n", w.ID, job.Name, r)
			job.UpdateStatus(Failed)
		}
		w.Perform()
	}()

	if job == nil {
		w.Channel <- fmt.Sprintf("job was nil for worker %d", w.ID)
		return
	}

	job.UpdateStatus(Processing)
	w.Channel <- fmt.Sprintf("job %s initiated by worker %d", job.Name, w.ID)
	err := w.Work(job.Data)
	if err != nil {
		job.UpdateStatus(Failed)
	} else {
		job.UpdateStatus(Completed)
	}
}
