package queue

import (
	"fmt"
	"strings"
)

type MiniQ[T any] struct {
	workers    []*Worker[T]
	queue      *Queue[T]
	logChannel chan string
	jobChannel chan string
}

func CreateMiniQ[T any](logChannel chan string) *MiniQ[T] {
	jobChannel := make(chan string)

	miniQ := MiniQ[T]{
		workers: []*Worker[T]{},
		queue: &Queue[T]{
			jobs: []*Job[T]{},
		},
		logChannel: logChannel,
		jobChannel: jobChannel,
	}

	go miniQ.Listen()
	return &miniQ
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
		go availableWorker.Perform(job)
	}
}

func (wp *MiniQ[T]) Listen() {
	for {
		msg := <-wp.jobChannel
		fmt.Println("JOB:", msg)
		if strings.Contains(msg, "WORKER_FREED") {
			wp.findFirstAvailableWorker()
		}
		if strings.Contains(msg, "JOB_ADDED") {
			wp.findFirstAvailableWorker()
		}
	}
}

func (wp *MiniQ[T]) AddJob(job *Job[T]) {
	wp.queue.enqueue(job)
	wp.jobChannel <- fmt.Sprintf("JOB_ADDED: %d", job.ID)
}

func (wp *MiniQ[T]) AddWorker(work func(T) error) {
	wp.workers = append(wp.workers, &Worker[T]{
		ID:         len(wp.workers),
		Work:       work,
		LogChannel: wp.logChannel,
		JobChannel: wp.jobChannel,
		Status:     Idle,
	})
}
