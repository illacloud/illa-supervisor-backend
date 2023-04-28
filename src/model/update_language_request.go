package model

type UpdateLanguageRequest struct {
	Language string `json:"language" validate:"oneof=en-US zh-CN ja-JP ko-KR cs-CZ da-DK de-DE el-GR es-ES fi-FI fr-FR it-IT nl-NL no-NO pl-PL pt-PT ru-RU ro-RO sv-SE uk-UA"`
}

func NewUpdateLanguageRequest() *UpdateLanguageRequest {
	return &UpdateLanguageRequest{}
}
