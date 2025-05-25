package handlers

import (
	"net/http"

	"github.com/AryaJayadi/MedTrace_api/internal/auth"
	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/drug"
	"github.com/AryaJayadi/MedTrace_api/internal/services"
	"github.com/labstack/echo/v4"
)

// DrugHandler handles HTTP requests for drugs
type DrugHandler struct {
	Service *services.DrugService
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
// @Failure 400 {object} response.BaseResponse "Invalid request payload or missing required fields"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /drugs [post]
// @Security BearerAuth
func (h *DrugHandler) CreateDrug(c echo.Context) error {
	var req drug.CreateDrugRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid request payload: " + err.Error()}})
	}

	if req.OwnerID == "" || req.BatchID == "" || req.DrugID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "OwnerID, BatchID, and DrugID are required"}})
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler CreateDrug: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.CreateDrug(contract, c.Request().Context(), &req)
	status := http.StatusCreated
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
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
// @Failure 400 {object} response.BaseResponse "Invalid Drug ID"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 404 {object} response.BaseResponse "Drug not found"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /drugs/{drugID} [get]
// @Security BearerAuth
func (h *DrugHandler) GetDrug(c echo.Context) error {
	drugID := c.Param("drugID")
	if drugID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Drug ID parameter is required"}})
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetDrug: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.GetDrug(contract, c.Request().Context(), drugID)
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
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /drugs/my [get]
// @Security BearerAuth
func (h *DrugHandler) GetMyDrugs(c echo.Context) error {
	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetMyDrugs: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.GetMyDrugs(contract, c.Request().Context())
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
// @Failure 400 {object} response.BaseResponse "Invalid Batch ID"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /drugs/batch/{batchID} [get]
// @Security BearerAuth
func (h *DrugHandler) GetDrugByBatch(c echo.Context) error {
	batchID := c.Param("batchID")
	if batchID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Batch ID parameter is required"}})
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetDrugByBatch: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.GetDrugByBatch(contract, c.Request().Context(), batchID)
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}

// GetMyAvailDrugs godoc
// @Summary Get available drugs owned by the caller
// @Description Retrieve all drug assets owned by the transaction submitter that are not currently in a pending transfer.
// @Tags drugs
// @Produce json
// @Success 200 {object} response.BaseListResponse[entity.Drug]
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /drugs/my/available [get]
// @Security BearerAuth
func (h *DrugHandler) GetMyAvailDrugs(c echo.Context) error {
	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetMyAvailDrugs: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.GetMyAvailDrugs(contract, c.Request().Context())
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}
