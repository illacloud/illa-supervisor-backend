package model

type VerificationRequest struct {
	Email string `json:"email" validate:"required"`
	Usage string `json:"usage" validate:"oneof=signup forgetpwd"`
}

func NewVerificationRequest() *VerificationRequest {
	return &VerificationRequest{}
}
