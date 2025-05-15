package drug

// CreateDrugRequest defines the structure for creating a new drug via API
type CreateDrugRequest struct {
	OwnerID string `json:"ownerID"`
	BatchID string `json:"batchID"`
	DrugID  string `json:"drugID"` // Client might suggest an ID or chaincode generates it.
	// Based on chaincode, it seems client provides it.
}
