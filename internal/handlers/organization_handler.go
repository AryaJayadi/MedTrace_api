package handlers

import (
	"net/http"

	"github.com/AryaJayadi/MedTrace_api/internal/auth"
	"github.com/AryaJayadi/MedTrace_api/internal/services"
	"github.com/labstack/echo/v4"
)

type OrganizationHandler struct {
	Service *services.OrganizationService
}

func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{Service: service}
}

// GetOrganizationByID godoc
// @Summary Get an organization by ID
// @Description Retrieve a specific organization from the ledger
// @Tags organizations
// @Produce json
// @Param id path string true "Organization ID"
// @Success 200 {object} response.BaseValueResponse[entity.Organization]
// @Failure 400 {object} response.BaseResponse "Invalid organization ID"
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 404 {object} response.BaseResponse "Organization not found"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /organizations/{id} [get]
// @Security BearerAuth
func (h *OrganizationHandler) GetOrganizationByID(c echo.Context) error {
	orgID := c.Param("id")
	if orgID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Organization ID parameter is required"}})
	}

	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetOrganizationByID: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.GetOrganizationByID(contract, c.Request().Context(), orgID)
	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}
	return c.JSON(status, resp)
}

// GetOrganizations godoc
// @Summary Get all organizations
// @Description Retrieve all organizations from the ledger
// @Tags organizations
// @Produce json
// @Success 200 {object} response.BaseListResponse[entity.Organization]
// @Failure 401 {object} response.BaseResponse "Unauthorized - JWT invalid or missing"
// @Failure 500 {object} response.BaseResponse "Internal server error or Fabric error"
// @Router /organizations [get]
// @Security BearerAuth
func (h *OrganizationHandler) GetOrganizations(c echo.Context) error {
	contract, err := auth.GetContractFromContext(c)
	if err != nil {
		c.Logger().Errorf("Handler GetOrganizations: Failed to get contract from context: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to access network resources"}})
	}

	resp := h.Service.GetOrganizations(contract, c.Request().Context())

	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}

	return c.JSON(status, resp)
}
