package services

import (
	"context"
	"encoding/json"

	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/transfer"
	"github.com/AryaJayadi/MedTrace_api/internal/models/entity"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// TransferService handles transfer-related operations
type TransferService struct {
	contract *client.Contract
}

// NewTransferService creates a new TransferService
func NewTransferService(contract *client.Contract) *TransferService {
	return &TransferService{contract: contract}
}

// CreateTransfer calls the CreateTransfer chaincode function
func (s *TransferService) CreateTransfer(ctx context.Context, req *transfer.CreateTransferRequest) response.BaseValueResponse[entity.Transfer] {
	// Directly marshal the API DTO. Its JSON tags are now set to match chaincode expectations.
	ccReqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to marshal CreateTransfer request: %v", err)
	}

	resultBytes, err := s.contract.SubmitTransaction("CreateTransfer", string(ccReqJSON))
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to submit CreateTransfer transaction: %v", err)
	}

	var transferEntity entity.Transfer
	err = json.Unmarshal(resultBytes, &transferEntity)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to unmarshal CreateTransfer result: %v", err)
	}
	return response.SuccessValueResponse(transferEntity)
}

// GetTransfer calls the GetTransfer chaincode function
func (s *TransferService) GetTransfer(ctx context.Context, transferID string) response.BaseValueResponse[entity.Transfer] {
	resultBytes, err := s.contract.EvaluateTransaction("GetTransfer", transferID)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to evaluate GetTransfer transaction: %v", err)
	}
	if len(resultBytes) == 0 {
		return response.ErrorValueResponse[entity.Transfer](404, "Transfer %s not found", transferID)
	}

	var transferEntity entity.Transfer
	err = json.Unmarshal(resultBytes, &transferEntity)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to unmarshal GetTransfer result: %v", err)
	}
	return response.SuccessValueResponse(transferEntity)
}

// generic helper for GetMy...Transfer functions
func (s *TransferService) getMyTransfersByType(ctx context.Context, chaincodeFunc string) response.BaseListResponse[entity.Transfer] {
	resultBytes, err := s.contract.EvaluateTransaction(chaincodeFunc)
	if err != nil {
		return response.ErrorListResponse[entity.Transfer](500, "Failed to evaluate %s transaction: %v", chaincodeFunc, err)
	}

	var transfers []entity.Transfer
	err = json.Unmarshal(resultBytes, &transfers)
	if err != nil {
		return response.ErrorListResponse[entity.Transfer](500, "Failed to unmarshal %s result: %v", chaincodeFunc, err)
	}

	transfersPtrs := make([]*entity.Transfer, len(transfers))
	for i := range transfers {
		transfersPtrs[i] = &transfers[i]
	}
	return response.SuccessListResponse(transfersPtrs)
}

// GetMyOutTransfer calls the GetMyOutTransfer chaincode function
func (s *TransferService) GetMyOutTransfer(ctx context.Context) response.BaseListResponse[entity.Transfer] {
	return s.getMyTransfersByType(ctx, "GetMyOutTransfer")
}

// GetMyInTransfer calls the GetMyInTransfer chaincode function
func (s *TransferService) GetMyInTransfer(ctx context.Context) response.BaseListResponse[entity.Transfer] {
	return s.getMyTransfersByType(ctx, "GetMyInTransfer")
}

// GetMyTransfers calls the GetMyTransfers chaincode function (all for the user)
func (s *TransferService) GetMyTransfers(ctx context.Context) response.BaseListResponse[entity.Transfer] {
	return s.getMyTransfersByType(ctx, "GetMyTransfers")
}

// AcceptTransfer calls the AcceptTransfer chaincode function
func (s *TransferService) AcceptTransfer(ctx context.Context, req *transfer.ProcessTransferRequest) response.BaseValueResponse[entity.Transfer] {
	// Directly marshal the API DTO. Its JSON tags are now set to match chaincode expectations.
	ccReqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to marshal AcceptTransfer request: %v", err)
	}

	resultBytes, err := s.contract.SubmitTransaction("AcceptTransfer", string(ccReqJSON))
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to submit AcceptTransfer transaction: %v", err)
	}

	var transferEntity entity.Transfer
	err = json.Unmarshal(resultBytes, &transferEntity)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to unmarshal AcceptTransfer result: %v", err)
	}
	return response.SuccessValueResponse(transferEntity)
}

// RejectTransfer calls the RejectTransfer chaincode function
func (s *TransferService) RejectTransfer(ctx context.Context, req *transfer.ProcessTransferRequest) response.BaseValueResponse[entity.Transfer] {
	// Directly marshal the API DTO. Its JSON tags are now set to match chaincode expectations.
	ccReqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to marshal RejectTransfer request: %v", err)
	}

	resultBytes, err := s.contract.SubmitTransaction("RejectTransfer", string(ccReqJSON))
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to submit RejectTransfer transaction: %v", err)
	}

	var transferEntity entity.Transfer
	err = json.Unmarshal(resultBytes, &transferEntity)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to unmarshal RejectTransfer result: %v", err)
	}
	return response.SuccessValueResponse(transferEntity)
}
