package queue

import (
	"errors"

	"github.com/fs1g17/MiniQ/internal/store"
)

var errNoJobInQueue = errors.New("queue is empty")

type MiniQ struct {
	jobStore *store.JobStore
	queue    *Queue
}

func CreateMiniQ(jobStore *store.JobStore) *MiniQ {
	miniQ := MiniQ{
		jobStore: jobStore,
		queue: &Queue{
			jobs: []*store.Job{},
		},
	}

	return &miniQ
}

func (wp *MiniQ) AddJob(data *store.AnyData) error {
	job := store.Job{
		Data: *data,
	}
	err := wp.jobStore.InsertJob(&job)
	if err != nil {
		return err
	}
	wp.queue.enqueue(&job)
	return nil
}

func (wp *MiniQ) GetJob() (*store.Job, error) {
	job := wp.queue.dequeue()
	if job == nil {
		return nil, errNoJobInQueue
	}

	wp.jobStore.UpdateJobStatus(job.ID, store.Processing)

	return job, nil
}
