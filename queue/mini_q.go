package queue

import (
	"fmt"
	"strings"
)

type MiniQ[T any] struct {
	workers []*Worker[T]
	queue   *Queue[T]
	channel chan string
}

func CreateMiniQ[T any](channel chan string) *MiniQ[T] {
	miniQ := MiniQ[T]{
		workers: []*Worker[T]{},
		queue: &Queue[T]{
			jobs: []*Job[T]{},
		},
		channel: channel,
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
		msg := <-wp.channel
		fmt.Println("HERE", msg)
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
	wp.channel <- fmt.Sprintf("JOB_ADDED: %s", job.Name)
}

func (wp *MiniQ[T]) AddWorker(work func(T) error) {
	wp.workers = append(wp.workers, &Worker[T]{
		ID:      len(wp.workers),
		Work:    work,
		Channel: wp.channel,
		Status:  Idle,
	})
}
