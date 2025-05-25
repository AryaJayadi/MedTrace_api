package handlers

import (
	"net/http"
	"time"

	"github.com/AryaJayadi/MedTrace_api/internal/auth"
	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/transfer"
	"github.com/AryaJayadi/MedTrace_api/internal/models/entity"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/AryaJayadi/MedTrace_api/internal/services"
	"github.com/labstack/echo/v4"
)

// TransferHandler handles HTTP requests for transfers
type TransferHandler struct {
	Service *services.TransferService
}

// NewTransferHandler creates a new TransferHandler
func NewTransferHandler(service *services.TransferService) *TransferHandler {
	return &TransferHandler{Service: service}
}

// sendResponse is a generic helper for Value responses
func (h *TransferHandler) sendValueResponse(c echo.Context, successStatus int, resp *response.BaseValueResponse[any]) error {
	if resp.Success {
		return c.JSON(successStatus, resp)
	}
	httpStatus := http.StatusInternalServerError
	if resp.Error != nil && resp.Error.Code != 0 {
		httpStatus = resp.Error.Code
	}
	return c.JSON(httpStatus, resp)
}

// sendListResponse is a generic helper for List responses
func (h *TransferHandler) sendListResponse(c echo.Context, successStatus int, resp *response.BaseListResponse[any]) error {
	if resp.Success {
		return c.JSON(successStatus, resp)
	}
	httpStatus := http.StatusInternalServerError
	if resp.Error != nil && resp.Error.Code != 0 {
		httpStatus = resp.Error.Code
	}
	return c.JSON(httpStatus, resp)
}

// CreateTransfer godoc
// @Summary Create a new transfer
// @Description Initiate a new transfer of drugs
// @Tags transfers
// @Accept json
// @Produce json
// @Param transfer body transfer.CreateTransferRequest true "Transfer details. DrugsID and ReceiverID are required."
// @Success 201 {object} response.BaseValueResponse[entity.Transfer]
// @Failure 400 {object} response.BaseResponse "Invalid request payload or missing required fields"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /transfers [post]
// @Security BearerAuth
func (h *TransferHandler) CreateTransfer(c echo.Context) error {
	var req transfer.CreateTransferRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusBadRequest, Message: "Invalid request payload: " + err.Error()}})
	}
	if req.ReceiverID == "" || len(req.DrugsID) == 0 {
		return c.JSON(http.StatusBadRequest, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusBadRequest, Message: "ReceiverID and at least one DrugID are required"}})
	}
	if req.TransferDate == nil {
		now := time.Now()
		req.TransferDate = &now
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler CreateTransfer: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusInternalServerError, Message: "Failed to access network resources"}})
	}

	resp := h.Service.CreateTransfer(contract, c.Request().Context(), &req)
	if resp.Success {
		return c.JSON(http.StatusCreated, resp)
	}
	httpStatus := http.StatusInternalServerError
	if resp.Error != nil && resp.Error.Code != 0 {
		httpStatus = resp.Error.Code
	}
	return c.JSON(httpStatus, resp)
}

// GetTransfer godoc
// @Summary Get a transfer by ID
// @Description Retrieve a specific transfer from the ledger
// @Tags transfers
// @Produce json
// @Param id path string true "Transfer ID"
// @Success 200 {object} response.BaseValueResponse[entity.Transfer]
// @Failure 400 {object} response.BaseResponse "Invalid Transfer ID"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 404 {object} response.BaseResponse "Transfer not found"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /transfers/{id} [get]
// @Security BearerAuth
func (h *TransferHandler) GetTransfer(c echo.Context) error {
	transferID := c.Param("id")
	if transferID == "" {
		return c.JSON(http.StatusBadRequest, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusBadRequest, Message: "Transfer ID parameter is required"}})
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetTransfer: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusInternalServerError, Message: "Failed to access network resources"}})
	}

	resp := h.Service.GetTransfer(contract, c.Request().Context(), transferID)
	if resp.Success {
		return c.JSON(http.StatusOK, resp)
	}
	httpStatus := http.StatusInternalServerError
	if resp.Error != nil && resp.Error.Code != 0 {
		httpStatus = resp.Error.Code
	}
	return c.JSON(httpStatus, resp)
}

// GetMyOutTransfer godoc
// @Summary Get outgoing transfers for the caller
// @Description Retrieve all outgoing transfers initiated by the transaction submitter
// @Tags transfers
// @Produce json
// @Success 200 {object} response.BaseListResponse[entity.Transfer]
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /transfers/my/outgoing [get]
// @Security BearerAuth
func (h *TransferHandler) GetMyOutTransfer(c echo.Context) error {
	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetMyOutTransfer: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, response.BaseListResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusInternalServerError, Message: "Failed to access network resources"}})
	}

	resp := h.Service.GetMyOutTransfer(contract, c.Request().Context())
	if resp.Success {
		return c.JSON(http.StatusOK, resp)
	}
	httpStatus := http.StatusInternalServerError
	if resp.Error != nil && resp.Error.Code != 0 {
		httpStatus = resp.Error.Code
	}
	return c.JSON(httpStatus, resp)
}

// GetMyInTransfer godoc
// @Summary Get incoming transfers for the caller
// @Description Retrieve all incoming transfers destined for the transaction submitter
// @Tags transfers
// @Produce json
// @Success 200 {object} response.BaseListResponse[entity.Transfer]
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /transfers/my/incoming [get]
// @Security BearerAuth
func (h *TransferHandler) GetMyInTransfer(c echo.Context) error {
	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetMyInTransfer: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, response.BaseListResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusInternalServerError, Message: "Failed to access network resources"}})
	}

	resp := h.Service.GetMyInTransfer(contract, c.Request().Context())
	if resp.Success {
		return c.JSON(http.StatusOK, resp)
	}
	httpStatus := http.StatusInternalServerError
	if resp.Error != nil && resp.Error.Code != 0 {
		httpStatus = resp.Error.Code
	}
	return c.JSON(httpStatus, resp)
}

// GetMyTransfers godoc
// @Summary Get all (incoming and outgoing) transfers for the caller
// @Description Retrieve all transfers associated with the transaction submitter
// @Tags transfers
// @Produce json
// @Success 200 {object} response.BaseListResponse[entity.Transfer]
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /transfers/my [get]
// @Security BearerAuth
func (h *TransferHandler) GetMyTransfers(c echo.Context) error {
	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetMyTransfers: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, response.BaseListResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusInternalServerError, Message: "Failed to access network resources"}})
	}

	resp := h.Service.GetMyTransfers(contract, c.Request().Context())
	if resp.Success {
		return c.JSON(http.StatusOK, resp)
	}
	httpStatus := http.StatusInternalServerError
	if resp.Error != nil && resp.Error.Code != 0 {
		httpStatus = resp.Error.Code
	}
	return c.JSON(httpStatus, resp)
}

// AcceptTransfer godoc
// @Summary Accept an incoming transfer
// @Description Mark an incoming transfer as accepted
// @Tags transfers
// @Accept json
// @Produce json
// @Param transfer body transfer.ProcessTransferRequest true "Transfer acceptance details. TransferID and ReceiveDate are required."
// @Success 200 {object} response.BaseValueResponse[entity.Transfer]
// @Failure 400 {object} response.BaseResponse "Invalid request payload or missing TransferID"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /transfers/accept [post]
// @Security BearerAuth
func (h *TransferHandler) AcceptTransfer(c echo.Context) error {
	var req transfer.ProcessTransferRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusBadRequest, Message: "Invalid request payload: " + err.Error()}})
	}
	if req.TransferID == "" {
		return c.JSON(http.StatusBadRequest, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusBadRequest, Message: "TransferID is required"}})
	}
	if req.ReceiveDate == nil {
		now := time.Now()
		req.ReceiveDate = &now
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler AcceptTransfer: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusInternalServerError, Message: "Failed to access network resources"}})
	}

	resp := h.Service.AcceptTransfer(contract, c.Request().Context(), &req)
	if resp.Success {
		return c.JSON(http.StatusOK, resp)
	}
	httpStatus := http.StatusInternalServerError
	if resp.Error != nil && resp.Error.Code != 0 {
		httpStatus = resp.Error.Code
	}
	return c.JSON(httpStatus, resp)
}

// RejectTransfer godoc
// @Summary Reject an incoming transfer
// @Description Mark an incoming transfer as rejected
// @Tags transfers
// @Accept json
// @Produce json
// @Param transfer body transfer.ProcessTransferRequest true "Transfer rejection details. Only TransferID is required."
// @Success 200 {object} response.BaseValueResponse[entity.Transfer]
// @Failure 400 {object} response.BaseResponse "Invalid request payload or missing TransferID"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /transfers/reject [post]
// @Security BearerAuth
func (h *TransferHandler) RejectTransfer(c echo.Context) error {
	var req transfer.ProcessTransferRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusBadRequest, Message: "Invalid request payload: " + err.Error()}})
	}
	if req.TransferID == "" {
		return c.JSON(http.StatusBadRequest, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusBadRequest, Message: "TransferID is required"}})
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler RejectTransfer: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, response.BaseValueResponse[entity.Transfer]{Success: false, Error: &response.ErrorInfo{Code: http.StatusInternalServerError, Message: "Failed to access network resources"}})
	}

	resp := h.Service.RejectTransfer(contract, c.Request().Context(), &req)
	if resp.Success {
		return c.JSON(http.StatusOK, resp)
	}
	httpStatus := http.StatusInternalServerError
	if resp.Error != nil && resp.Error.Code != 0 {
		httpStatus = resp.Error.Code
	}
	return c.JSON(httpStatus, resp)
}
