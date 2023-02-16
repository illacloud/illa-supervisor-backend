package model

type UpdateLanguageRequest struct {
	Language string `json:"language" validate:"oneof=zh-CN en-US ja-JP ko-KR"`
}

func NewUpdateLanguageRequest() *UpdateLanguageRequest {
	return &UpdateLanguageRequest{}
}
