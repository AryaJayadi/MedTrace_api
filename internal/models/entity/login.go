package entity

type LoginResponseData struct {
	Token   string `json:"Token"`
	OrgID   string `json:"OrgId"`
	Message string `json:"Message"`
}
