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

func (s *BatchService) CreateBatch(ctx context.Context, req *batch.BatchCreate) response.BaseValueResponse[entity.Batch] {
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
