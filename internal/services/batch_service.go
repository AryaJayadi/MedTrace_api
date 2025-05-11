package services

import (
	"context"
	"encoding/json"
	"fmt"

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

func (s *BatchService) CreateBatch(ctx context.Context, req *batch.BatchCreate) response.BaseValueResponse {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse(400, "Failed to marshal request")
	}

	resp, err := s.contract.SubmitTransaction("CreateBatch", string(reqJSON))
	if err != nil {
		msg := fmt.Sprintf("Failed to submit transaction to Fabric: %v", err)
		return response.ErrorValueResponse(500, msg)
	}

	var batch entity.Batch
	if err := json.Unmarshal(resp, &batch); err != nil {
		return response.ErrorValueResponse(500, "Failed to unmarshal Fabric response")
	}

	return response.SuccessValueResponse(batch)
}
