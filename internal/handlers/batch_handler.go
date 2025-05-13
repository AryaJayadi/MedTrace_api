package handlers

import (
	"net/http"

	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/batch"
	"github.com/AryaJayadi/MedTrace_api/internal/services"
	"github.com/labstack/echo/v4"
)

type BatchHandler struct {
	service *services.BatchService
}

func NewBatchHandler(service *services.BatchService) *BatchHandler {
	return &BatchHandler{service: service}
}

func (h *BatchHandler) CreateBatch(c echo.Context) error {
	var req batch.CreateBatch
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

func (h *BatchHandler) GetAllBatches(c echo.Context) error {
	resp := h.service.GetAllBatches(c.Request().Context())

	status := http.StatusOK
	if !resp.Success {
		status = http.StatusInternalServerError
	}

	return c.JSON(status, resp)
}

func (h *BatchHandler) UpdateBatch(c echo.Context) error {
	var req batch.UpdateBatch
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	resp := h.service.UpdateBatch(c.Request().Context(), &req)

	status := http.StatusOK
	if !resp.Success {
		status = http.StatusInternalServerError
	}

	return c.JSON(status, resp)
}
