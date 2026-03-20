package queue

import (
	"fmt"
	"time"
)

type Worker[T any] struct {
	ID      int
	Work    func(T) error
	Queue   *Queue[T]
	Channel chan string
}

func (w *Worker[T]) Perform() {
	for {
		time.Sleep(200 * time.Millisecond)
		job := w.Queue.Dequeue()
		if job == nil {
			w.Channel <- fmt.Sprintf("job was nil for worker %d", w.ID)
			continue
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
}
