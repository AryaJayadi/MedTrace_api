package handlers

import (
	"net/http"

	"github.com/AryaJayadi/MedTrace_api/internal/services"
	"github.com/labstack/echo/v4"
)

type OrganizationHandler struct {
	service *services.OrganizationService
}

func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

func (h *OrganizationHandler) GetOrganizations(c echo.Context) error {
	resp := h.service.GetOrganizations(c.Request().Context())

	status := http.StatusOK
	if !resp.Success {
		status = http.StatusInternalServerError
	}

	return c.JSON(status, resp)
}
