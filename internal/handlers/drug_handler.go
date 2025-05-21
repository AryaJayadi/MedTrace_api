package handlers

import (
	"net/http"

	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/drug"
	"github.com/AryaJayadi/MedTrace_api/internal/services"
	"github.com/labstack/echo/v4"
)

// DrugHandler handles HTTP requests for drugs
type DrugHandler struct {
	Service *services.DrugService // Renamed from 'service' to 'Service' for convention
}

// NewDrugHandler creates a new DrugHandler
func NewDrugHandler(drugService *services.DrugService) *DrugHandler {
	return &DrugHandler{Service: drugService}
}

// CreateDrug godoc
// @Summary Create a new drug
// @Description Create a new drug asset on the ledger. The ID is returned by the chaincode.
// @Tags drugs
// @Accept json
// @Produce json
// @Param drug body drug.CreateDrugRequest true "Drug to create. OwnerID, BatchID, DrugID are required."
// @Success 201 {object} response.BaseValueResponse[string]
// @Failure 400 {object} response.BaseResponse "{ \"error\": \"Bad Request\" }"
// @Failure 500 {object} response.BaseResponse "{ \"error\": \"Internal Server Error\" }"
// @Router /drugs [post]
func (h *DrugHandler) CreateDrug(c echo.Context) error {
	var req drug.CreateDrugRequest
	if err := c.Bind(&req); err != nil {
		// Using the generic ErrorResponse from your response package structure for consistency if available
		// Assuming BaseResponse has a structure for error, or use map[string]string
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid request payload: " + err.Error()}})
	}

	if req.OwnerID == "" || req.BatchID == "" || req.DrugID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "OwnerID, BatchID, and DrugID are required"}})
	}

	resp := h.Service.CreateDrug(c.Request().Context(), &req)
	status := http.StatusCreated
	if !resp.Success {
		status = resp.Error.Code // Assuming ErrorInfo has a Code field for HTTP status
		if status == 0 {         // Default to 500 if code not set or invalid
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}

// GetDrug godoc
// @Summary Get a drug by ID
// @Description Retrieve a specific drug asset from the ledger
// @Tags drugs
// @Produce json
// @Param drugID path string true "Drug ID"
// @Success 200 {object} response.BaseValueResponse[entity.Drug]
// @Failure 404 {object} response.BaseResponse "{ \"error\": \"Drug not found\" }"
// @Failure 500 {object} response.BaseResponse "{ \"error\": \"Internal Server Error\" }"
// @Router /drugs/{drugID} [get]
func (h *DrugHandler) GetDrug(c echo.Context) error {
	drugID := c.Param("drugID")
	if drugID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Drug ID parameter is required"}})
	}

	resp := h.Service.GetDrug(c.Request().Context(), drugID)
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}

// GetMyDrugs godoc
// @Summary Get drugs owned by the caller
// @Description Retrieve all drug assets owned by the transaction submitter from the ledger
// @Tags drugs
// @Produce json
// @Success 200 {object} response.BaseListResponse[entity.Drug]
// @Failure 500 {object} response.BaseResponse "{ \"error\": \"Internal Server Error\" }"
// @Router /drugs/my [get]
func (h *DrugHandler) GetMyDrugs(c echo.Context) error {
	resp := h.Service.GetMyDrugs(c.Request().Context())
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}

// GetDrugByBatch godoc
// @Summary Get drugs by batch ID
// @Description Retrieve all drug assets associated with a specific batch ID from the ledger
// @Tags drugs
// @Produce json
// @Param batchID path string true "Batch ID"
// @Success 200 {object} response.BaseListResponse[entity.Drug]
// @Failure 400 {object} response.BaseResponse "{ \"error\": \"Bad Request\" }"
// @Failure 500 {object} response.BaseResponse "{ \"error\": \"Internal Server Error\" }"
// @Router /drugs/batch/{batchID} [get]
func (h *DrugHandler) GetDrugByBatch(c echo.Context) error {
	batchID := c.Param("batchID")
	if batchID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Batch ID parameter is required"}})
	}

	resp := h.Service.GetDrugByBatch(c.Request().Context(), batchID)
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}
