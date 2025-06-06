package entity

// Drug entity based on chaincode model
type Drug struct {
	ID            string `json:"ID"`
	BatchID       string `json:"BatchID"`
	OwnerID       string `json:"OwnerID"`
	Location      string `json:"Location"`
	IsTransferred bool   `json:"isTransferred"`
	TransferID    string `json:"TransferID,omitempty"` // omitempty if it might not always be present
}
