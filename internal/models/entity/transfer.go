package entity

import (
	"time"

	"github.com/AryaJayadi/MedTrace_api/internal/utils"
)

// Transfer entity based on chaincode model (updated to value types)
type Transfer struct {
	ID           string             `json:"ID"`
	IsAccepted   bool               `json:"isAccepted"`
	ReceiveDate  utils.OptionalTime `json:"ReceiveDate"`
	ReceiverID   string             `json:"ReceiverID"`
	SenderID     string             `json:"SenderID"`
	TransferDate time.Time          `json:"TransferDate"`
}
