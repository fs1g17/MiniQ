package api

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/fs1g17/MiniQ/internal/store"
	"github.com/labstack/echo/v5"
)

type MiniQueue interface {
	AddJob(data *store.AnyData) (*store.Job, error)
	CompleteJob(jobID int, success bool) error
	GetJob() (*store.Job, error)
	GetJobs() []*store.Job
	AssignJob(jobID int) error
}

type QueueHandler struct {
	miniq   MiniQueue
	mu      sync.RWMutex
	clients map[chan struct{}]struct{}
	jobs    []store.Job
}

func NewQueueHandler(miniq MiniQueue) *QueueHandler {
	return &QueueHandler{
		miniq:   miniq,
		clients: make(map[chan struct{}]struct{}),
		jobs:    make([]store.Job, 0),
	}
}

type addJobRequest struct {
	Data store.AnyData `json:"data"`
}

func (r *addJobRequest) validate() error {
	if r.Data == nil {
		return errors.New("data misisng")
	}
	return nil
}

func (h *QueueHandler) HandleAddJob(c *echo.Context) error {
	var req addJobRequest
	err := c.Bind(&req)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := req.validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	_, err = h.miniq.AddJob(&req.Data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

outer:
	for clientPingChan := range h.clients {
		select {
		case clientPingChan <- struct{}{}:
			break outer // sends to first available channel
		default:
			// CLient is slow or not listening, skip
		}
	}

	return c.JSON(http.StatusNoContent, nil)
}

type completeJobRequest struct {
	JobID   int  `json:"jobID"`
	Success bool `json:"success"`
}

func (h *QueueHandler) HandleCompleteJob(c *echo.Context) error {
	var req completeJobRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := h.miniq.CompleteJob(req.JobID, req.Success)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *QueueHandler) HandleGetJob(c *echo.Context) error {
	job, err := h.miniq.GetJob()

	if err != nil {
		return c.JSON(http.StatusOK, map[string]any{"job": nil})
	}

	return c.JSON(http.StatusOK, map[string]any{"job": job})
}

func (h *QueueHandler) HandleGetJobs(c *echo.Context) error {
	jobs := h.miniq.GetJobs()
	return c.JSON(http.StatusOK, map[string]any{"jobs": jobs})
}

func (h *QueueHandler) HandlePollJob(c *echo.Context) error {
	pingChan := make(chan struct{}, 1)

	h.mu.Lock()
	h.clients[pingChan] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, pingChan)
		h.mu.Unlock()
		close(pingChan)
	}()

	ctx := c.Request().Context()

	// if jobs available, get directly from the queue
	job, _ := h.miniq.GetJob()
	if job != nil {
		c.Response().Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusOK, map[string]any{"job": job})
		return nil
	}

	// if no job available, wait inside the long poll until one is added
	select {
	case <-pingChan:
		job, _ := h.miniq.GetJob()
		if job == nil {
			c.JSON(http.StatusNoContent, nil)
			return nil
		}
		c.Response().Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusOK, map[string]any{"job": job})
	case <-time.After(30 * time.Second):
		c.JSON(http.StatusNoContent, nil)
	case <-ctx.Done():
		// just return nil
	}

	return nil
}
