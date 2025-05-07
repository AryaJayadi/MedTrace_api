package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

type BatchHandler struct {
	service *services.BatchService
}

func NewBatchHandler(service *services.BatchService) *BatchHandler {
	return &BatchHandler{service: service}
}

func (h *BatchHandler) CreateBatch(c echo.Context) error {
	var req batch.BatchCreate
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	resp := h.service.CreateBatch(c.Request().Context(), &req)

	status := http.StatusCreated
	if !resp.Success {
		status = http.StatusInternalServerError
	}

	return c.JSON(status, resp)
}
