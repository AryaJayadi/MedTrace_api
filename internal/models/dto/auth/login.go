package auth

type LoginResponseData struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
	Message      string `json:"Message"`
	OrgID        string `json:"OrgId"`
}
