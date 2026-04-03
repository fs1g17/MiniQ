package queue

import (
	"errors"
	"fmt"

	"github.com/fs1g17/MiniQ/internal/store"
)

var errNoJobInQueue = errors.New("queue is empty")

type MiniQ struct {
	jobStore *store.JobStore
	queue    *Queue
}

func CreateMiniQ(jobStore *store.JobStore) *MiniQ {
	jobs, err := jobStore.GetQueuedJobs()
	if err != nil {
		panic(fmt.Sprintf("failed to recreated job queue from db %v", err))
	}

	miniQ := MiniQ{
		jobStore: jobStore,
		queue: &Queue{
			jobs: jobs,
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

func (wp *MiniQ) GetJobs() []*store.Job {
	return wp.queue.getJobs()
}
