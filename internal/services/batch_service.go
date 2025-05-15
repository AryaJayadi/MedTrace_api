package services

import (
	"context"
	"encoding/json"

	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/batch"
	"github.com/AryaJayadi/MedTrace_api/internal/models/entity"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type BatchService struct {
	contract *client.Contract
}

func NewBatchService(contract *client.Contract) *BatchService {
	return &BatchService{contract: contract}
}

func (s *BatchService) CreateBatch(ctx context.Context, req *batch.CreateBatch) response.BaseValueResponse[entity.Batch] {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to marshal request: %v", err)
	}

	resp, err := s.contract.SubmitTransaction("CreateBatch", string(reqJSON))
	if err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to submit transaction to Fabric: %v", err)
	}

	var batch entity.Batch
	if err := json.Unmarshal(resp, &batch); err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to unmarshal Fabric response: %v", err)
	}

	return response.SuccessValueResponse(batch)
}

func (s *BatchService) GetAllBatches(ctx context.Context) response.BaseListResponse[entity.Batch] {
	resp, err := s.contract.EvaluateTransaction("GetAllBatches")
	if err != nil {
		return response.ErrorListResponse[entity.Batch](500, "Failed to evaluate transaction: %v", err)
	}

	var batches []entity.Batch
	if err := json.Unmarshal(resp, &batches); err != nil {
		return response.ErrorListResponse[entity.Batch](500, "Failed to unmarshal Fabric response: %v", err)
	}

	ptrList := make([]*entity.Batch, len(batches))
	for i := range batches {
		ptrList[i] = &batches[i]
	}

	return response.SuccessListResponse(ptrList)
}

func (s *BatchService) UpdateBatch(ctx context.Context, req *batch.UpdateBatch) response.BaseValueResponse[entity.Batch] {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to marshal request: %v", err)
	}

	resp, err := s.contract.SubmitTransaction("UpdateBatch", string(reqJSON))
	if err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to submit transaction to Fabric: %v", err)
	}

	var batch entity.Batch
	if err := json.Unmarshal(resp, &batch); err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to unmarshal Fabric response: %v", err)
	}

	return response.SuccessValueResponse(batch)
}

// GetBatchByID calls the GetBatch chaincode function
func (s *BatchService) GetBatchByID(ctx context.Context, batchID string) response.BaseValueResponse[entity.Batch] {
	resultBytes, err := s.contract.EvaluateTransaction("GetBatch", batchID)
	if err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to evaluate GetBatch transaction: %v", err)
	}
	if len(resultBytes) == 0 {
		return response.ErrorValueResponse[entity.Batch](404, "Batch %s not found", batchID)
	}

	var batchEntity entity.Batch
	err = json.Unmarshal(resultBytes, &batchEntity)
	if err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to unmarshal batch data: %v", err)
	}
	return response.SuccessValueResponse(batchEntity)
}

// BatchExists calls the BatchExists chaincode function
func (s *BatchService) BatchExists(ctx context.Context, batchID string) response.BaseValueResponse[bool] {
	resultBytes, err := s.contract.EvaluateTransaction("BatchExists", batchID)
	if err != nil {
		return response.ErrorValueResponse[bool](500, "Failed to evaluate BatchExists transaction: %v", err)
	}
	// The chaincode BatchExists returns a boolean marshalled as JSON string (e.g., "true" or "false")
	var exists bool
	err = json.Unmarshal(resultBytes, &exists)
	if err != nil {
		// It's possible resultBytes is not a valid JSON bool, e.g. empty if record doesn't exist and CC returns nothing
		// Or if chaincode returns non-JSON string. For safety, treat unmarshal error as not found or error.
		return response.ErrorValueResponse[bool](500, "Failed to unmarshal BatchExists result: %v. Raw: %s", err, string(resultBytes))
	}
	return response.SuccessValueResponse(exists)
}
