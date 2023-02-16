package model

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required"`
}

func NewChangePasswordRequest() *ChangePasswordRequest {
	return &ChangePasswordRequest{}
}
