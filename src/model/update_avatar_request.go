package model

type UpdateAvatarRequest struct {
	Avatar string `json:"avatar" validate:"required"`
}

func NewUpdateAvatarRequest() *UpdateAvatarRequest {
	return &UpdateAvatarRequest{}
}
