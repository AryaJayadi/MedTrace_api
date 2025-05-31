package auth

type PayloadRefreshToken struct {
	RefreshToken string `json:"RefreshToken" validate:"required"`
}
