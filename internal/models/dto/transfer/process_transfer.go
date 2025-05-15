package transfer

import "time"

// ProcessTransferRequest defines the structure for accepting or rejecting a transfer
// Its JSON tags must match the fields expected by the chaincode's ProcessTransfer DTO
type ProcessTransferRequest struct {
	TransferID  string     `json:"transferID"`            // Chaincode expects "transferID".
	ReceiveDate *time.Time `json:"ReceiveDate,omitempty"` // Chaincode expects "ReceiveDate".
}
