package entity

import (
	"encoding/json"
	"time"
)

type HistoryDrug struct {
	Drug      *Drug     `json:"Drug"`
	TxID      string    `json:"TxID"`
	Timestamp time.Time `json:"Timestamp"`
	IsDelete  bool      `json:"IsDelete"`
}

type historyDrugInput struct {
	Record    *Drug     `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

func (h *HistoryDrug) UnmarshalJSON(data []byte) error {
	var input historyDrugInput
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	h.Drug = input.Record
	h.TxID = input.TxId
	h.Timestamp = input.Timestamp
	h.IsDelete = input.IsDelete
	return nil
}
