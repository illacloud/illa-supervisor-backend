package model

type ForgetPasswordRequest struct {
	Email             string `json:"email" validate:"required"`
	NewPassword       string `json:"newPassword" validate:"required"`
	VerificationCode  string `json:"verificationCode" validate:"required"`
	VerificationToken string `json:"verificationToken" validate:"required"`
}

func NewForgetPasswordRequest() *ForgetPasswordRequest {
	return &ForgetPasswordRequest{}
}
