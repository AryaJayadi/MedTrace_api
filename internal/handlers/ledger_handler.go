package handlers

import (
	"net/http"

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
// @Success 200 {object} response.BaseValueResponse[string]
// @Failure 500 {object} response.BaseResponse
// @Router /ledger/init [post]
func (h *LedgerHandler) InitLedger(c echo.Context) error {
	resp := h.Service.InitLedger(c.Request().Context())
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}
