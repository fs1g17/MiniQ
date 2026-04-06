package api

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fs1g17/MiniQ/internal/queue"
	"github.com/fs1g17/MiniQ/internal/store"
	"github.com/labstack/echo/v5"
)

type QueueHandler struct {
	miniq   *queue.MiniQ
	mu      sync.RWMutex
	clients map[chan store.Job]struct{}
	jobs    []store.Job
}

func NewQueueHandler(miniq *queue.MiniQ) *QueueHandler {
	return &QueueHandler{
		miniq:   miniq,
		clients: make(map[chan store.Job]struct{}),
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
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := req.validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	job, err := h.miniq.AddJob(&req.Data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

outer:
	for clientChan := range h.clients {
		select {
		case clientChan <- *job:
			break outer // sends to first available channel
		default:
			// CLient is slow or not listening, skip
		}
	}

	return nil
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
	jobChan := make(chan store.Job, 1)

	h.mu.Lock()
	h.clients[jobChan] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, jobChan)
		h.mu.Unlock()
		close(jobChan)
	}()

	ctx := c.Request().Context()

	// if jobs available, get directly from the queue
	job, _ := h.miniq.GetJob()
	log.Printf("HERE, job is: %v", job)
	if job != nil {
		c.Response().Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusOK, map[string]any{"job": job})
		return nil
	}

	// if no job available, wait inside the long poll until one is added
	select {
	case job := <-jobChan:
		h.miniq.AssignJob(job.ID)
		c.Response().Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusOK, map[string]any{"job": job})
	case <-time.After(30 * time.Second):
		c.JSON(http.StatusNoContent, nil)
	case <-ctx.Done():
		return nil
	}

	return nil
}
