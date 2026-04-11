package queue

import (
	"errors"
	"fmt"

	"github.com/fs1g17/MiniQ/internal/store"
)

var errNoJobInQueue = errors.New("queue is empty")

type JobStore interface {
	GetQueuedJobs() ([]*store.Job, error)
	InsertJob(job *store.Job) error
	UpdateJobStatus(jobId int, jobStatus store.JobStatus) error
}

type MiniQ struct {
	jobStore JobStore
	queue    *Queue
}

func CreateMiniQ(jobStore JobStore) *MiniQ {
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

func (wp *MiniQ) AddJob(data *store.AnyData) (*store.Job, error) {
	job := store.Job{
		Data: *data,
	}
	err := wp.jobStore.InsertJob(&job)
	if err != nil {
		return nil, err
	}
	wp.queue.enqueue(&job)
	return &job, nil
}

func (wp *MiniQ) GetJob() (*store.Job, error) {
	job := wp.queue.dequeue()
	if job == nil {
		return nil, errNoJobInQueue
	}

	wp.jobStore.UpdateJobStatus(job.ID, store.Processing)

	return job, nil
}

func (wp *MiniQ) AssignJob(jobID int) error {
	err := wp.jobStore.UpdateJobStatus(jobID, store.Processing)
	return err
}

func (wp *MiniQ) CompleteJob(jobID int, success bool) error {
	var jobStatus store.JobStatus
	if success {
		jobStatus = store.Completed
	} else {
		jobStatus = store.Failed
	}
	err := wp.jobStore.UpdateJobStatus(jobID, jobStatus)
	return err
}

func (wp *MiniQ) GetJobs() []*store.Job {
	return wp.queue.getJobs()
}
