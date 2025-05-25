package handlers

import (
	"net/http"

	"github.com/AryaJayadi/MedTrace_api/internal/auth"
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
// @Failure 400 {object} response.BaseResponse "Invalid request payload"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /batches [post]
// @Security BearerAuth
func (h *BatchHandler) CreateBatch(c echo.Context) error {
	var req batch.CreateBatch
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid request payload: " + err.Error()}})
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler CreateBatch: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.CreateBatch(contract, c.Request().Context(), &req)
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
// @Failure 400 {object} response.BaseResponse "Invalid Batch ID"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 404 {object} response.BaseResponse "Batch not found"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /batches/{id} [get]
// @Security BearerAuth
func (h *BatchHandler) GetBatchByID(c echo.Context) error {
	batchID := c.Param("id")
	if batchID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Batch ID parameter is required"}})
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetBatchByID: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.GetBatchByID(contract, c.Request().Context(), batchID)
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
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /batches [get]
// @Security BearerAuth
func (h *BatchHandler) GetAllBatches(c echo.Context) error {
	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetAllBatches: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.GetAllBatches(contract, c.Request().Context())
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
// @Failure 400 {object} response.BaseResponse "Invalid request payload"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 404 {object} response.BaseResponse "Batch not found to update"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /batches [patch]  // Consider /batches/{id} if updating a specific batch by ID in path
// @Security BearerAuth
func (h *BatchHandler) UpdateBatch(c echo.Context) error {
	var req batch.UpdateBatch
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid request payload: " + err.Error()}})
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler UpdateBatch: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.UpdateBatch(contract, c.Request().Context(), &req)
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
// @Failure 400 {object} response.BaseResponse "Invalid Batch ID"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /batches/{id}/exists [get]
// @Security BearerAuth
func (h *BatchHandler) BatchExists(c echo.Context) error {
	batchID := c.Param("id")
	if batchID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Batch ID parameter is required"}})
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler BatchExists: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.BatchExists(contract, c.Request().Context(), batchID)
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}
