package model

type UpdateIsTutorialViewedRequest struct {
	IsTutorialViewed bool `json:"isTutorialViewed" validate:"boolean"`
}

func NewUpdateIsTutorialViewedRequest() *UpdateIsTutorialViewedRequest {
	return &UpdateIsTutorialViewedRequest{}
}
