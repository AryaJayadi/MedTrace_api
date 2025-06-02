package services

import (
	"context"
	"encoding/json"

	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/drug"
	"github.com/AryaJayadi/MedTrace_api/internal/models/entity"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// DrugService handles drug-related operations.
// It no longer stores the contract directly.
type DrugService struct {
	// No contract field here
}

// NewDrugService creates a new DrugService.
// It no longer takes a contract as a parameter.
func NewDrugService() *DrugService {
	return &DrugService{}
}

// CreateDrug calls the CreateDrug chaincode function using the provided contract.
func (s *DrugService) CreateDrug(contract *client.Contract, ctx context.Context, req *drug.CreateDrugRequest) response.BaseValueResponse[string] {
	// Chaincode CreateDrug returns drugID string, not the full drug object directly from that call.
	resultBytes, err := contract.SubmitTransaction("CreateDrug", req.OwnerID, req.BatchID, req.DrugID)
	if err != nil {
		return response.ErrorValueResponse[string](500, "Failed to submit CreateDrug transaction: %v", err)
	}
	// The result from chaincode is the drugID string
	drugID := string(resultBytes)
	return response.SuccessValueResponse(drugID)
}

// GetDrug calls the GetDrug chaincode function using the provided contract.
func (s *DrugService) GetDrug(contract *client.Contract, ctx context.Context, drugID string) response.BaseValueResponse[entity.Drug] {
	resultBytes, err := contract.EvaluateTransaction("GetDrug", drugID)
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

// GetMyDrugs calls the GetMyDrug chaincode function using the provided contract.
func (s *DrugService) GetMyDrugs(contract *client.Contract, ctx context.Context) response.BaseListResponse[entity.Drug] {
	resultBytes, err := contract.EvaluateTransaction("GetMyDrug")
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

// GetDrugByBatch calls the GetDrugByBatch chaincode function using the provided contract.
func (s *DrugService) GetDrugByBatch(contract *client.Contract, ctx context.Context, batchID string) response.BaseListResponse[entity.Drug] {
	resultBytes, err := contract.EvaluateTransaction("GetDrugByBatch", batchID)
	if err != nil {
		return response.ErrorListResponse[entity.Drug](500, "Failed to evaluate GetDrugByBatch transaction: %v", err)
	}
	if len(resultBytes) == 0 {
		// Return empty list if no drugs found for the batch, not necessarily an error
		return response.SuccessListResponse([]*entity.Drug{})
	}

	var drugs []entity.Drug
	err = json.Unmarshal(resultBytes, &drugs)
	if err != nil {
		return response.ErrorListResponse[entity.Drug](500, "Failed to unmarshal drugs data for GetDrugByBatch: %v", err)
	}

	drugsPtrs := make([]*entity.Drug, len(drugs))
	for i := range drugs {
		drugsPtrs[i] = &drugs[i]
	}

	return response.SuccessListResponse(drugsPtrs)
}

// GetMyAvailDrugs calls the GetMyAvailDrugs chaincode function using the provided contract.
func (s *DrugService) GetMyAvailDrugs(contract *client.Contract, ctx context.Context) response.BaseListResponse[entity.Drug] {
	resultBytes, err := contract.EvaluateTransaction("GetMyAvailDrugs")
	if err != nil {
		return response.ErrorListResponse[entity.Drug](500, "Failed to evaluate GetMyAvailDrugs transaction: %v", err)
	}

	var drugs []entity.Drug
	err = json.Unmarshal(resultBytes, &drugs)
	if err != nil {
		return response.ErrorListResponse[entity.Drug](500, "Failed to unmarshal drugs data for GetMyAvailDrugs: %v", err)
	}

	drugsPtrs := make([]*entity.Drug, len(drugs))
	for i := range drugs {
		drugsPtrs[i] = &drugs[i]
	}

	return response.SuccessListResponse(drugsPtrs)
}

func (s *DrugService) GetHistoryDrug(contract *client.Contract, ctx context.Context, drugID string) response.BaseListResponse[entity.HistoryDrug] {
	resultBytes, err := contract.EvaluateTransaction("GetHistoryDrug", drugID)
	if err != nil {
		return response.ErrorListResponse[entity.HistoryDrug](500, "Failed to evaluate GetHistoryDrug transaction: %v", err)
	}

	var records []entity.HistoryDrug
	err = json.Unmarshal(resultBytes, &records)
	if err != nil {
		return response.ErrorListResponse[entity.HistoryDrug](500, "Failed to unmarshal history drug data for GetHistoryDrug: %v", err)
	}

	recordsPtrs := make([]*entity.HistoryDrug, len(records))
	for i := range records {
		recordsPtrs[i] = &records[i]
	}

	return response.SuccessListResponse(recordsPtrs)
}
