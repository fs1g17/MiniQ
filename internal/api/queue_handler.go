package api

import (
	"errors"
	"net/http"

	"github.com/fs1g17/MiniQ/internal/queue"
	"github.com/fs1g17/MiniQ/internal/store"
	"github.com/labstack/echo/v5"
)

type QueueHandler struct {
	miniq *queue.MiniQ
}

func NewQueueHandler(miniq *queue.MiniQ) *QueueHandler {
	return &QueueHandler{
		miniq: miniq,
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

	err := h.miniq.AddJob(&req.Data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
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
