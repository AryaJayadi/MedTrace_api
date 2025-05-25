package handlers

import (
	"net/http"

	"github.com/AryaJayadi/MedTrace_api/internal/auth"
	"github.com/AryaJayadi/MedTrace_api/internal/services"
	"github.com/labstack/echo/v4"
)

// LedgerHandler handles HTTP requests for ledger operations
type LedgerHandler struct {
	Service *services.LedgerService
}

// NewLedgerHandler creates a new LedgerHandler
func NewLedgerHandler(service *services.LedgerService) *LedgerHandler {
	return &LedgerHandler{Service: service}
}

// InitLedger godoc
// @Summary Initialize the ledger
// @Description Run the InitLedger chaincode function. This is typically a one-time setup.
// @Tags ledger
// @Produce json
// @Success 200 {object} response.BaseValueResponse[string] "Successfully initialized ledger"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /ledger/init [post]
// @Security BearerAuth
func (h *LedgerHandler) InitLedger(c echo.Context) error {
	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler InitLedger: Failed to get contract from context: %v", err)
		// Assuming response.BaseValueResponse structure for error consistency
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.InitLedger(contract, c.Request().Context())
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}
