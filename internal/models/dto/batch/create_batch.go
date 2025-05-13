package batch

import (
	"time"
)

type CreateBatch struct {
	Amount         int       `json:"Amount" xml:"Amount" form:"Amount"`                         // Amount of drugs in the batch
	Description    string    `json:"Description" xml:"Description" form:"Description"`          // Drug description
	Manufacturer   string    `json:"Manufacturer" xml:"Manufacturer" form:"Manufacturer"`       // Manufacturer name
	Name           string    `json:"Name" xml:"Name" form:"Name"`                               // Drug name
	ExpiryDate     time.Time `json:"ExpiryDate" xml:"ExpiryDate" form:"ExpiryDate"`             // Drug expiry date
	Location       string    `json:"Location" xml:"Location" form:"Location"`                   // Location of the drug
	ProductionDate time.Time `json:"ProductionDate" xml:"ProductionDate" form:"ProductionDate"` // Drug production date
	Status         string    `json:"Status" xml:"Status" form:"Status"`                         // Drug status
}
