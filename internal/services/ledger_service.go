package services

import (
	"context"

	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// LedgerService handles ledger-related operations.
// It no longer stores the contract directly.
type LedgerService struct {
	// No contract field here
}

// NewLedgerService creates a new LedgerService.
// It no longer takes a contract as a parameter.
func NewLedgerService() *LedgerService {
	return &LedgerService{}
}

// InitLedger calls the InitLedger chaincode function using the provided contract.
// The chaincode InitLedger function doesn't return a specific value on success, just an error if it fails.
// So, we'll return a simple success message.
func (s *LedgerService) InitLedger(contract *client.Contract, ctx context.Context) response.BaseValueResponse[string] {
	// ctx is available if needed for future use (e.g. timeouts, cancellation), but not directly used by SubmitTransaction here.
	_, err := contract.SubmitTransaction("InitLedger") // Result not typically used for InitLedger
	if err != nil {
		return response.ErrorValueResponse[string](500, "Failed to submit InitLedger transaction: %v", err)
	}
	return response.SuccessValueResponse("Ledger initialized successfully.")
}
