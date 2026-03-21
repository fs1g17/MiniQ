package queue

type MiniQ[T any] struct {
	workers []*Worker[T]
	queue   *Queue[T]
}

func CreateMiniQ[T any]() *MiniQ[T] {
	return &MiniQ[T]{
		workers: []*Worker[T]{},
		queue: &Queue[T]{
			jobs: []*Job[T]{},
		},
	}
}

func (wp *MiniQ[T]) findFirstAvailableWorker() {
	job := wp.queue.dequeue()
	if job == nil {
		return // no jobs
	}

	var availableWorker *Worker[T] = nil
	for _, worker := range wp.workers {
		if workerStatus := worker.GetStatus(); workerStatus == Busy {
			continue
		}
		availableWorker = worker
		break
	}
	if availableWorker != nil {
		availableWorker.SetStatus(Busy)
		go availableWorker.Perform(job, wp.findFirstAvailableWorker)
	}
}

func (wp *MiniQ[T]) AddJob(job *Job[T]) {
	wp.queue.enqueue(job)
	wp.findFirstAvailableWorker()
}

func (wp *MiniQ[T]) AddWorker(work func(T) error, channel chan string) {
	wp.workers = append(wp.workers, &Worker[T]{
		ID:      len(wp.workers),
		Work:    work,
		Channel: channel,
		Status:  Idle,
	})
}
