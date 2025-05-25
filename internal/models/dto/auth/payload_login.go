package auth

type PayloadLogin struct {
	Organization string `json:"organization" validate:"required"`
	Password     string `json:"password" validate:"required"`
}
