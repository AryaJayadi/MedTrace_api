package services

import (
	"context"
	"encoding/json"

	"github.com/AryaJayadi/MedTrace_api/internal/models/entity"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// OrganizationService handles business logic for organizations.
// It no longer stores the contract directly.
type OrganizationService struct {
	// No contract field here
}

// NewOrganizationService creates a new OrganizationService.
// It no longer takes a contract as a parameter.
func NewOrganizationService() *OrganizationService {
	return &OrganizationService{}
}

// GetOrganizationByID retrieves a specific organization from the ledger using the provided contract.
func (s *OrganizationService) GetOrganizationByID(contract *client.Contract, ctx context.Context, orgID string) response.BaseValueResponse[entity.Organization] {
	resultBytes, err := contract.EvaluateTransaction("GetOrganization", orgID)
	if err != nil {
		return response.ErrorValueResponse[entity.Organization](500, "Failed to evaluate GetOrganization transaction: %v", err)
	}
	if len(resultBytes) == 0 {
		return response.ErrorValueResponse[entity.Organization](404, "Organization %s not found", orgID)
	}

	var org entity.Organization
	err = json.Unmarshal(resultBytes, &org)
	if err != nil {
		return response.ErrorValueResponse[entity.Organization](500, "Failed to unmarshal organization data: %v", err)
	}
	return response.SuccessValueResponse(org)
}

// GetOrganizations retrieves all organizations from the ledger using the provided contract.
func (s *OrganizationService) GetOrganizations(contract *client.Contract, ctx context.Context) response.BaseListResponse[entity.Organization] {
	resp, err := contract.EvaluateTransaction("GetAllOrganizations")
	if err != nil {
		return response.ErrorListResponse[entity.Organization](500, "Failed to evaluate transaction to Fabric: %v", err)
	}

	var organizations []entity.Organization
	if err := json.Unmarshal(resp, &organizations); err != nil {
		return response.ErrorListResponse[entity.Organization](500, "Failed to unmarshal Fabric response: %v", err)
	}

	ptrList := make([]*entity.Organization, len(organizations))
	for i := range organizations {
		ptrList[i] = &organizations[i]
	}

	return response.SuccessListResponse(ptrList)
}
