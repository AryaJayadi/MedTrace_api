package services

import (
	"context"

	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// LedgerService handles ledger-related operations
type LedgerService struct {
	contract *client.Contract
}

// NewLedgerService creates a new LedgerService
func NewLedgerService(contract *client.Contract) *LedgerService {
	return &LedgerService{contract: contract}
}

// InitLedger calls the InitLedger chaincode function
// The chaincode InitLedger function doesn't return a specific value on success, just an error if it fails.
// So, we'll return a simple success message.
func (s *LedgerService) InitLedger(ctx context.Context) response.BaseValueResponse[string] {
	_, err := s.contract.SubmitTransaction("InitLedger") // Result not typically used for InitLedger
	if err != nil {
		return response.ErrorValueResponse[string](500, "Failed to submit InitLedger transaction: %v", err)
	}
	return response.SuccessValueResponse("Ledger initialized successfully.")
}
