package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/fs1g17/MiniQ/internal/queue"
	"github.com/fs1g17/MiniQ/internal/store"
	"github.com/labstack/echo/v5"
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

type mockMiniQ struct {
	miniQ   *queue.MiniQ
	readyCh chan struct{}
}

func (m *mockMiniQ) AddJob(data *store.AnyData) (*store.Job, error) {
	return m.miniQ.AddJob(data)
}

func (m *mockMiniQ) CompleteJob(jobID int, success bool) error {
	return m.miniQ.CompleteJob(jobID, success)
}

func (m *mockMiniQ) GetJob() (*store.Job, error) {
	defer func() {
		m.readyCh <- struct{}{}
	}()

	return m.miniQ.GetJob()
}

func (m *mockMiniQ) GetJobs() []*store.Job {
	return m.miniQ.GetJobs()
}

func (m *mockMiniQ) AssignJob(jobID int) error {
	return m.miniQ.AssignJob(jobID)
}

func TestEdgeCase(t *testing.T) {
	jobStore := NewMockJobStore()
	miniQ := queue.CreateMiniQ(jobStore)

	mockQ := mockMiniQ{
		miniQ:   miniQ,
		readyCh: make(chan struct{}),
	}
	queueHandler := NewQueueHandler(&mockQ)

	req := httptest.NewRequest(http.MethodGet, "/pollJob", nil)
	rec := httptest.NewRecorder()
	ctx := echo.NewContext(req, rec)

	go func() {
		queueHandler.HandlePollJob(ctx)
	}()

	<-mockQ.readyCh

	queueHandler.miniq.AddJob(&store.AnyData{"message": "hello world"})
	time.Sleep(1 * time.Second)
	job, err := miniQ.GetJob()

	assert.Nil(t, job, "Job should be nil")
	assert.Error(t, err, "Expected error to be empty queue")
}
