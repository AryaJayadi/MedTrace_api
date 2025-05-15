package handlers

import (
	"net/http"

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
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /organizations/{id} [get]
func (h *OrganizationHandler) GetOrganizationByID(c echo.Context) error {
	orgID := c.Param("id")
	if orgID == "" {
		// Consider using response.ErrorValueResponse for consistency if you have it for bad requests
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]interface{}{"code": http.StatusBadRequest, "message": "Organization ID parameter is required"}})
	}

	resp := h.Service.GetOrganizationByID(c.Request().Context(), orgID)
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
// @Failure 500 {object} response.BaseResponse
// @Router /organizations [get]
func (h *OrganizationHandler) GetOrganizations(c echo.Context) error {
	resp := h.Service.GetOrganizations(c.Request().Context())

	status := http.StatusOK
	if !resp.Success {
		status = resp.Error.Code
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}

	return c.JSON(status, resp)
}
