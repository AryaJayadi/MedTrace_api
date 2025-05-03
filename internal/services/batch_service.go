package services

import (
	"context"
	"encoding/json"

	"github.com/AryaJayadi/SupplyChain_api/internal/models/dto/batch"
	"github.com/AryaJayadi/SupplyChain_api/internal/models/entity"
	"github.com/AryaJayadi/SupplyChain_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type BatchService struct {
	contract *client.Contract
}

func NewBatchService(contract *client.Contract) *BatchService {
	return &BatchService{contract: contract}
}

func (s *BatchService) CreateBatch(ctx context.Context, req *batch.BatchCreate) response.BaseValueResponse {
	resp, err := s.contract.SubmitTransaction("CreateBatch", req)
	if err != nil {
		return response.ErrorValueResponse(500, "Failed to submit transaction")
	}

	var batch entity.Batch
	err = json.Unmarshal(resp, &batch)
	if err != nil {
		return response.ErrorValueResponse(500, "Error unmarshaling JSON")
	}

	return response.SuccessValueResponse(batch)
}
