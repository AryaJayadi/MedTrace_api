package services

import (
	"context"
	"encoding/json"

	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/transfer"
	"github.com/AryaJayadi/MedTrace_api/internal/models/entity"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// TransferService handles transfer-related operations.
// It no longer stores the contract directly.
type TransferService struct {
	// No contract field here
}

// NewTransferService creates a new TransferService.
// It no longer takes a contract as a parameter.
func NewTransferService() *TransferService {
	return &TransferService{}
}

// CreateTransfer calls the CreateTransfer chaincode function using the provided contract.
func (s *TransferService) CreateTransfer(contract *client.Contract, ctx context.Context, req *transfer.CreateTransferRequest) response.BaseValueResponse[entity.Transfer] {
	ccReqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to marshal CreateTransfer request: %v", err)
	}

	resultBytes, err := contract.SubmitTransaction("CreateTransfer", string(ccReqJSON))
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

// GetTransfer calls the GetTransfer chaincode function using the provided contract.
func (s *TransferService) GetTransfer(contract *client.Contract, ctx context.Context, transferID string) response.BaseValueResponse[entity.Transfer] {
	resultBytes, err := contract.EvaluateTransaction("GetTransfer", transferID)
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

// getMyTransfersByType is a generic helper for GetMy...Transfer functions using the provided contract.
func (s *TransferService) getMyTransfersByType(contract *client.Contract, ctx context.Context, chaincodeFunc string) response.BaseListResponse[entity.Transfer] {
	resultBytes, err := contract.EvaluateTransaction(chaincodeFunc)
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

// GetMyOutTransfer calls the GetMyOutTransfer chaincode function using the provided contract.
func (s *TransferService) GetMyOutTransfer(contract *client.Contract, ctx context.Context) response.BaseListResponse[entity.Transfer] {
	return s.getMyTransfersByType(contract, ctx, "GetMyOutTransfer")
}

// GetMyInTransfer calls the GetMyInTransfer chaincode function using the provided contract.
func (s *TransferService) GetMyInTransfer(contract *client.Contract, ctx context.Context) response.BaseListResponse[entity.Transfer] {
	return s.getMyTransfersByType(contract, ctx, "GetMyInTransfer")
}

// GetMyTransfers calls the GetMyTransfers chaincode function (all for the user) using the provided contract.
func (s *TransferService) GetMyTransfers(contract *client.Contract, ctx context.Context) response.BaseListResponse[entity.Transfer] {
	return s.getMyTransfersByType(contract, ctx, "GetMyTransfers")
}

// AcceptTransfer calls the AcceptTransfer chaincode function using the provided contract.
func (s *TransferService) AcceptTransfer(contract *client.Contract, ctx context.Context, req *transfer.ProcessTransferRequest) response.BaseValueResponse[entity.Transfer] {
	ccReqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to marshal AcceptTransfer request: %v", err)
	}

	resultBytes, err := contract.SubmitTransaction("AcceptTransfer", string(ccReqJSON))
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

// RejectTransfer calls the RejectTransfer chaincode function using the provided contract.
func (s *TransferService) RejectTransfer(contract *client.Contract, ctx context.Context, req *transfer.ProcessTransferRequest) response.BaseValueResponse[entity.Transfer] {
	ccReqJSON, err := json.Marshal(req)
	if err != nil {
		return response.ErrorValueResponse[entity.Transfer](500, "Failed to marshal RejectTransfer request: %v", err)
	}

	resultBytes, err := contract.SubmitTransaction("RejectTransfer", string(ccReqJSON))
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
