package queue

import (
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"

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

func TestCompleteJob(t *testing.T) {
	tests := []struct {
		name    string
		success bool
		want    store.JobStatus
		message string
	}{
		{
			name:    "success",
			success: true,
			want:    store.Completed,
			message: "Expected job to be in completed state",
		},
		{
			name:    "failure",
			success: false,
			want:    store.Failed,
			message: "Expected job to be in failed state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobStore := NewMockJobStore()
			jobStore.jobs = append(jobStore.jobs, &store.Job{
				ID:        1,
				Status:    store.Processing,
				Data:      store.AnyData{"message": "hello world!"},
				Attempts:  0,
				CreatedAt: time.Now(),
			})
			miniq := CreateMiniQ(jobStore)

			miniq.CompleteJob(1, tt.success)

			assert.Equal(t, tt.want, jobStore.jobs[0].Status, tt.message)
		})
	}
}

func TestParallelJobs(t *testing.T) {
	jobStore := NewMockJobStore()
	jobStore.jobs = append(jobStore.jobs, &store.Job{
		ID:        1,
		Status:    store.Queued,
		Data:      store.AnyData{"message": "hello world!"},
		Attempts:  0,
		CreatedAt: time.Now(),
	})
	jobStore.jobs = append(jobStore.jobs, &store.Job{
		ID:        2,
		Status:    store.Queued,
		Data:      store.AnyData{"message": "hello world!"},
		Attempts:  0,
		CreatedAt: time.Now(),
	})
	miniq := CreateMiniQ(jobStore)

	jobChan := make(chan store.Job, 2)

	var wg sync.WaitGroup
	for range 2 {
		wg.Go(func() {
			job, err := miniq.GetJob()
			if err != nil {
				return
			}
			miniq.CompleteJob(job.ID, true)
			jobChan <- *job
		})
	}

	wg.Wait()
	close(jobChan)

	jobMap := make(map[int]struct{}, 2)
	jobSlice := make([]store.Job, 0, 2)
	for job := range jobChan {
		jobSlice = append(jobSlice, job)
		if _, ok := jobMap[job.ID]; ok == true {
			// value already exists, so some worker got same job
			t.Errorf("job with id %d given to multiple different workers", job.ID)
		}
		jobMap[job.ID] = struct{}{}
	}
	assert.Equal(t, 2, len(jobSlice), "Expected there to be 2 complete jobs")
}
