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
