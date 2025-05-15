package services

import (
	"context"
	"encoding/json"

	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/drug"
	"github.com/AryaJayadi/MedTrace_api/internal/models/entity"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// DrugService handles drug-related operations
type DrugService struct {
	contract *client.Contract
}

// NewDrugService creates a new DrugService
func NewDrugService(contract *client.Contract) *DrugService {
	return &DrugService{contract: contract}
}

// CreateDrug calls the CreateDrug chaincode function
func (s *DrugService) CreateDrug(ctx context.Context, req *drug.CreateDrugRequest) response.BaseValueResponse[string] {
	// Chaincode CreateDrug returns drugID string, not the full drug object directly from that call.
	resultBytes, err := s.contract.SubmitTransaction("CreateDrug", req.OwnerID, req.BatchID, req.DrugID)
	if err != nil {
		return response.ErrorValueResponse[string](500, "Failed to submit CreateDrug transaction: %v", err)
	}
	// The result from chaincode is the drugID string
	drugID := string(resultBytes)
	return response.SuccessValueResponse(drugID)
}

// GetDrug calls the GetDrug chaincode function
func (s *DrugService) GetDrug(ctx context.Context, drugID string) response.BaseValueResponse[entity.Drug] {
	resultBytes, err := s.contract.EvaluateTransaction("GetDrug", drugID)
	if err != nil {
		return response.ErrorValueResponse[entity.Drug](500, "Failed to evaluate GetDrug transaction: %v", err)
	}
	if len(resultBytes) == 0 { // Check if result is empty, indicating not found
		return response.ErrorValueResponse[entity.Drug](404, "Drug %s not found", drugID)
	}

	var drugEntity entity.Drug
	err = json.Unmarshal(resultBytes, &drugEntity)
	if err != nil {
		return response.ErrorValueResponse[entity.Drug](500, "Failed to unmarshal drug data for GetDrug: %v", err)
	}
	return response.SuccessValueResponse(drugEntity)
}

// GetMyDrugs calls the GetMyDrug chaincode function
func (s *DrugService) GetMyDrugs(ctx context.Context) response.BaseListResponse[entity.Drug] {
	resultBytes, err := s.contract.EvaluateTransaction("GetMyDrug")
	if err != nil {
		return response.ErrorListResponse[entity.Drug](500, "Failed to evaluate GetMyDrug transaction: %v", err)
	}

	var drugs []entity.Drug // Changed from []*entity.Drug to []entity.Drug for direct unmarshal
	err = json.Unmarshal(resultBytes, &drugs)
	if err != nil {
		return response.ErrorListResponse[entity.Drug](500, "Failed to unmarshal drugs data for GetMyDrug: %v", err)
	}

	// Convert []entity.Drug to []*entity.Drug for the response type
	drugsPtrs := make([]*entity.Drug, len(drugs))
	for i := range drugs {
		drugsPtrs[i] = &drugs[i]
	}

	return response.SuccessListResponse(drugsPtrs)
}
