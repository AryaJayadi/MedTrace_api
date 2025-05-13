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
