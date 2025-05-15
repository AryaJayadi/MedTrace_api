package handlers

import (
	"net/http"

	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/batch"
	"github.com/AryaJayadi/MedTrace_api/internal/services"
	"github.com/labstack/echo/v4"
)

type BatchHandler struct {
	Service *services.BatchService
}

func NewBatchHandler(service *services.BatchService) *BatchHandler {
	return &BatchHandler{Service: service}
}

// CreateBatch godoc
// @Summary Create a new batch
// @Description Create a new batch of drugs. Manufacturer details are derived from the caller.
// @Tags batches
// @Accept json
// @Produce json
// @Param batch body batch.CreateBatch true "Batch creation details"
// @Success 201 {object} response.BaseValueResponse[entity.Batch]
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /batches [post]
func (h *BatchHandler) CreateBatch(c echo.Context) error {
	var req batch.CreateBatch
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid request payload: " + err.Error()}})
	}

	resp := h.Service.CreateBatch(c.Request().Context(), &req)
	status := http.StatusCreated
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}

// GetBatchByID godoc
// @Summary Get a batch by ID
// @Description Retrieve a specific batch from the ledger
// @Tags batches
// @Produce json
// @Param id path string true "Batch ID"
// @Success 200 {object} response.BaseValueResponse[entity.Batch]
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /batches/{id} [get]
func (h *BatchHandler) GetBatchByID(c echo.Context) error {
	batchID := c.Param("id")
	if batchID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Batch ID parameter is required"}})
	}
	resp := h.Service.GetBatchByID(c.Request().Context(), batchID)
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}

// GetAllBatches godoc
// @Summary Get all batches
// @Description Retrieve all batches from the ledger.
// @Tags batches
// @Produce json
// @Success 200 {object} response.BaseListResponse[entity.Batch]
// @Failure 500 {object} response.BaseResponse
// @Router /batches [get]
func (h *BatchHandler) GetAllBatches(c echo.Context) error {
	resp := h.Service.GetAllBatches(c.Request().Context())
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}

// UpdateBatch godoc
// @Summary Update an existing batch
// @Description Update details of an existing batch.
// @Tags batches
// @Accept json
// @Produce json
// @Param batch body batch.UpdateBatch true "Batch update details"
// @Success 200 {object} response.BaseValueResponse[entity.Batch]
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /batches [patch]
func (h *BatchHandler) UpdateBatch(c echo.Context) error {
	var req batch.UpdateBatch
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid request payload: " + err.Error()}})
	}

	resp := h.Service.UpdateBatch(c.Request().Context(), &req)
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}

// BatchExists godoc
// @Summary Check if a batch exists
// @Description Check for the existence of a batch by its ID.
// @Tags batches
// @Produce json
// @Param id path string true "Batch ID"
// @Success 200 {object} response.BaseValueResponse[bool]
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /batches/{id}/exists [get]
func (h *BatchHandler) BatchExists(c echo.Context) error {
	batchID := c.Param("id")
	if batchID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Batch ID parameter is required"}})
	}
	resp := h.Service.BatchExists(c.Request().Context(), batchID)
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}
