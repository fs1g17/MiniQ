package queue

import (
	"fmt"
	"slices"
	"testing"

	"github.com/fs1g17/MiniQ/internal/store"
	"github.com/stretchr/testify/assert"
)

type mockJobStore struct {
	jobs   []*store.Job
	nextID int
}

func (m *mockJobStore) GetQueuedJobs() ([]*store.Job, error) {
	result := make([]*store.Job, 0)
	for _, job := range m.jobs {
		if job.Status == store.Queued {
			result = append(result, job)
		}
	}
	return result, nil
}

func (m *mockJobStore) InsertJob(job *store.Job) error {
	job.ID = m.nextID
	m.nextID++
	m.jobs = append(m.jobs, job)
	return nil
}

func (m *mockJobStore) UpdateJobStatus(jobId int, jobStatus store.JobStatus) error {
	idx := slices.IndexFunc(m.jobs, func(j *store.Job) bool { return j.ID == jobId })
	if idx == -1 {
		return fmt.Errorf("job with id %d not found", jobId)
	}

	m.jobs[idx].Status = jobStatus

	return nil
}

func NewMockJobStore() *mockJobStore {
	return &mockJobStore{
		jobs:   make([]*store.Job, 0),
		nextID: 1,
	}
}

func TestHappyPath(t *testing.T) {
	jobStore := NewMockJobStore()
	miniq := CreateMiniQ(jobStore)

	_, err := miniq.AddJob(&store.AnyData{"message": "Hello world!"})
	assert.NoError(t, err, "adding a job shouldn't return an error")
	length := len(miniq.queue.jobs)
	assert.Equal(t, 1, length, "Expected length of job queue to be 1")

	job, err := miniq.GetJob()
	assert.NoError(t, err, "Expected to not get error on miniq.GetJob()")
	assert.Equal(t, 0, len(miniq.queue.jobs), "Expected job queue to be empty")
	assert.Equal(t, 1, job.ID, "Expected first job to have ID 1")
	assert.Equal(t, store.Processing, job.Status, "Expected added job status to be processing")
	assert.Equal(t, "Hello world!", job.Data["message"], "Expected data message to be present")
	assert.Equal(t, store.Processing, jobStore.jobs[0].Status, "Expected job store job to be processing")
}

func TestEmptyQueue(t *testing.T) {
	jobStore := NewMockJobStore()
	miniq := CreateMiniQ(jobStore)

	_, err := miniq.GetJob()
	assert.EqualError(t, err, errNoJobInQueue.Error(), "Expected to get empty queue error")
}
