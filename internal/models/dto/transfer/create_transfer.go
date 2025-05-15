package transfer

import "time"

// CreateTransferRequest defines the structure for creating a new transfer via API
// Its JSON tags must match the fields expected by the chaincode's CreateTransfer DTO
type CreateTransferRequest struct {
	DrugsID      []string   `json:"DrugsID"`      // List of drug IDs. Chaincode expects "DrugsID".
	ReceiverID   string     `json:"ReceiverID"`   // Receiver ID. Chaincode expects "ReceiverID".
	TransferDate *time.Time `json:"TransferDate"` // Transfer date. Chaincode expects "TransferDate".
	// SenderID is omitted as it's determined by the chaincode from the caller's identity.
}
