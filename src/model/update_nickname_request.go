package model

type UpdateNicknameRequest struct {
	Nickname string `json:"nickname" validate:"required"`
}

func NewUpdateNicknameRequest() *UpdateNicknameRequest {
	return &UpdateNicknameRequest{}
}
