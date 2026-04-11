package queue

import (
	"fmt"
	"slices"

	"github.com/fs1g17/MiniQ/internal/store"
)

type mockJobStore struct {
	jobs   []*store.Job
	nextID int
}

func (m *mockJobStore) GetQueuedJobs() ([]*store.Job, error) {
	jobs := make([]*store.Job, 0)
	return jobs, nil
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
