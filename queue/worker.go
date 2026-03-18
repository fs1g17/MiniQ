package queue

import (
	"fmt"
	"time"
)

type Worker[T any] struct {
	ID      int
	Work    func(T)
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
		w.Channel <- fmt.Sprintf("job %s initiated by worker %d", job.Name, w.ID)
		go w.Work(job.Data)
	}
}
