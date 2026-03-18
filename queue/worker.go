package queue

import (
	"fmt"
	"time"
)

type Worker[T any] struct {
	Work  func(T)
	Queue *Queue[T]
}

func (w *Worker[T]) Perform() {
	for {
		time.Sleep(200 * time.Millisecond)
		job := w.Queue.Dequeue()
		if job == nil {
			// we do nothing
			fmt.Println("job was nil")
			continue
		}
		go w.Work(job.Data)
	}
}
