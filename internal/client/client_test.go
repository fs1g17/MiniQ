package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fs1g17/MiniQ/internal/store"
	"github.com/stretchr/testify/assert"
)

type TestServer struct {
	notifiedSuccess *bool
}

func (ts *TestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/pollJob" {
		w.Header().Set("Content-Type", "application/json")
		exampleJob := &store.Job{
			ID:        1,
			Status:    store.Processing,
			Data:      store.AnyData{"message": "hello world!"},
			Attempts:  0,
			CreatedAt: time.Now(),
		}
		err := json.NewEncoder(w).Encode(map[string]any{"job": &exampleJob})
		if err != nil {
			http.Error(w, "Failed to encode example job", http.StatusInternalServerError)
			return
		}
	}
	if r.URL.Path == "/completeJob" {
		var body struct {
			Success bool `json:"success"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		*ts.notifiedSuccess = body.Success
		w.WriteHeader(http.StatusOK)
	}
}

func TestTaskPanicRecover(t *testing.T) {
	notified := false
	ts := &TestServer{notifiedSuccess: &notified}
	testServer := httptest.NewServer(ts)
	defer testServer.Close()

	client := &http.Client{
		Timeout: 35 * time.Second,
	}

	task(testServer.URL, client, func(jobResponse JobResponse) {
		panic("boom")
	})

	assert.NotNil(t, ts.notifiedSuccess)
	assert.False(t, *ts.notifiedSuccess, "expected failure notification after panic")
}
