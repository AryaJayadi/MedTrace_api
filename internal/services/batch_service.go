package services

import (
	"context"
	"encoding/json"

	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/batch"
	"github.com/AryaJayadi/MedTrace_api/internal/models/entity"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// BatchService handles business logic for batches.
// It no longer stores the contract directly.
type BatchService struct {
	// No contract field here
}

// NewBatchService creates a new BatchService.
// It no longer takes a contract as a parameter.
func NewBatchService() *BatchService {
	return &BatchService{}
}

// CreateBatch creates a new batch on the ledger using the provided contract.
func (s *BatchService) CreateBatch(contract *client.Contract, ctx context.Context, req *batch.CreateBatch) response.BaseValueResponse[entity.Batch] {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to marshal request: %v", err)
	}

	resp, err := contract.SubmitTransaction("CreateBatch", string(reqJSON))
	if err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to submit transaction to Fabric: %v", err)
	}

	var batchEntity entity.Batch // Renamed to avoid conflict with package name
	if err := json.Unmarshal(resp, &batchEntity); err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to unmarshal Fabric response: %v", err)
	}

	return response.SuccessValueResponse(batchEntity)
}

// GetAllBatches retrieves all batches from the ledger using the provided contract.
func (s *BatchService) GetAllBatches(contract *client.Contract, ctx context.Context) response.BaseListResponse[entity.Batch] {
	resp, err := contract.EvaluateTransaction("GetAllBatches")
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

// UpdateBatch updates an existing batch on the ledger using the provided contract.
func (s *BatchService) UpdateBatch(contract *client.Contract, ctx context.Context, req *batch.UpdateBatch) response.BaseValueResponse[entity.Batch] {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to marshal request: %v", err)
	}

	resp, err := contract.SubmitTransaction("UpdateBatch", string(reqJSON))
	if err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to submit transaction to Fabric: %v", err)
	}

	var batchEntity entity.Batch // Renamed to avoid conflict
	if err := json.Unmarshal(resp, &batchEntity); err != nil {
		return response.ErrorValueResponse[entity.Batch](500, "Failed to unmarshal Fabric response: %v", err)
	}

	return response.SuccessValueResponse(batchEntity)
}

// GetBatchByID retrieves a specific batch by ID from the ledger using the provided contract.
func (s *BatchService) GetBatchByID(contract *client.Contract, ctx context.Context, batchID string) response.BaseValueResponse[entity.Batch] {
	resultBytes, err := contract.EvaluateTransaction("GetBatch", batchID)
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

// BatchExists checks if a batch exists on the ledger using the provided contract.
func (s *BatchService) BatchExists(contract *client.Contract, ctx context.Context, batchID string) response.BaseValueResponse[bool] {
	resultBytes, err := contract.EvaluateTransaction("BatchExists", batchID)
	if err != nil {
		return response.ErrorValueResponse[bool](500, "Failed to evaluate BatchExists transaction: %v", err)
	}
	var exists bool
	err = json.Unmarshal(resultBytes, &exists)
	if err != nil {
		return response.ErrorValueResponse[bool](500, "Failed to unmarshal BatchExists result: %v. Raw: %s", err, string(resultBytes))
	}
	return response.SuccessValueResponse(exists)
}
