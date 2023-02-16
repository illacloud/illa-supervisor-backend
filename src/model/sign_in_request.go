package model

type SignInRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewSignInRequest() *SignInRequest {
	return &SignInRequest{}
}
