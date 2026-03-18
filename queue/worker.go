package queue

type Worker[T any] struct {
	Work  func(T)
	Queue *Queue[T]
}

func (w *Worker[T]) Perform() {
	job := w.Queue.Dequeue()
	w.Work(job.Data)
}
