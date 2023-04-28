package model

type SignUpRequest struct {
	Nickname          string `json:"nickname" validate:"required"`
	Email             string `json:"email" validate:"required"`
	Password          string `json:"password" validate:"required"`
	Language          string `json:"language" validate:"oneof=en-US zh-CN ja-JP ko-KR cs-CZ da-DK de-DE el-GR es-ES fi-FI fr-FR it-IT nl-NL no-NO pl-PL pt-PT ru-RU ro-RO sv-SE uk-UA"`
	IsSubscribed      bool   `json:"isSubscribed"` // is subscribed are optional
	IsTutorialViewed  bool   `json:"isTutorialViewed"`
	VerificationCode  string `json:"verificationCode"`  // @todo: add validate:"required" when uset smtp server configureable
	VerificationToken string `json:"verificationToken"` // @todo: add validate:"required" when uset smtp server configureable
	InviteToken       string `json:"inviteToken"`       // invite token are optional
}

func NewSignUpRequest() *SignUpRequest {
	return &SignUpRequest{}
}

func (u *SignUpRequest) IsSignUpWithInviteLink() bool {
	if len(u.InviteToken) != 0 {
		return true
	}
	return false
}

func (u *SignUpRequest) ExportEmail() string {
	return u.Email
}

func (u *SignUpRequest) ExportInviteToken() string {
	return u.InviteToken
}
