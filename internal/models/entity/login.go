package entity

type LoginResponseData struct {
	Token   string `json:"token"`
	OrgID   string `json:"orgId"`
	Message string `json:"message"`
}
