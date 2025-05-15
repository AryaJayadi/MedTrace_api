package services

import (
	"context"
	"encoding/json"

	"github.com/AryaJayadi/MedTrace_api/internal/models/entity"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type OrganizationService struct {
	contract *client.Contract
}

func NewOrganizationService(contract *client.Contract) *OrganizationService {
	return &OrganizationService{contract: contract}
}

func (s *OrganizationService) GetOrganizationByID(ctx context.Context, orgID string) response.BaseValueResponse[entity.Organization] {
	resultBytes, err := s.contract.EvaluateTransaction("GetOrganization", orgID)
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

func (s *OrganizationService) GetOrganizations(ctx context.Context) response.BaseListResponse[entity.Organization] {
	resp, err := s.contract.EvaluateTransaction("GetAllOrganizations")
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
