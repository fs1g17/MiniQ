package queue

import (
	"fmt"
	"strings"

	"github.com/fs1g17/MiniQ/store"
)

type MiniQ struct {
	workers    []*Worker
	queue      *Queue
	logChannel chan string
	jobChannel chan string
}

func CreateMiniQ(logChannel chan string) *MiniQ {
	jobChannel := make(chan string)

	miniQ := MiniQ{
		workers: []*Worker{},
		queue: &Queue{
			jobs: []*store.Job{},
		},
		logChannel: logChannel,
		jobChannel: jobChannel,
	}

	go miniQ.Listen()
	return &miniQ
}

func (wp *MiniQ) findFirstAvailableWorker() {
	job := wp.queue.dequeue()
	if job == nil {
		return // no jobs
	}

	var availableWorker *Worker = nil
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

func (wp *MiniQ) Listen() {
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

func (wp *MiniQ) AddJob(job *store.Job) {
	wp.queue.enqueue(job)
	wp.jobChannel <- fmt.Sprintf("JOB_ADDED: %d", job.ID)
}

func (wp *MiniQ) AddWorker(work func(store.AnyData) error) {
	wp.workers = append(wp.workers, &Worker{
		ID:         len(wp.workers),
		Work:       work,
		LogChannel: wp.logChannel,
		JobChannel: wp.jobChannel,
		Status:     Idle,
	})
}
