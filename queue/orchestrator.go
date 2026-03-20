package queue

type Orchestrator[T any] struct {
	Queue   *Queue[T]
	Workers []*Worker[T]
}
