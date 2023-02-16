package model

type UpdateUsernameRequest struct {
	Nickname string `json:"nickname" validate:"required"`
}

func NewUpdateUsernameRequest() *UpdateUsernameRequest {
	return &UpdateUsernameRequest{}
}
