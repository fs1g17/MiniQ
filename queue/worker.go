package queue

type Worker[T any] struct {
	Work func(T)
}
