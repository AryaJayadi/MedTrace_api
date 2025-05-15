package entity

import "time"

// Transfer entity based on chaincode model (updated to value types)
type Transfer struct {
	ID           string    `json:"ID"`
	IsAccepted   bool      `json:"isAccepted"`
	ReceiveDate  time.Time `json:"ReceiveDate,omitempty"`
	ReceiverID   string    `json:"ReceiverID"`
	SenderID     string    `json:"SenderID"`
	TransferDate time.Time `json:"TransferDate,omitempty"`
}
