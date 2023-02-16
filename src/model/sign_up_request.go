package model

type SignUpRequest struct {
	Nickname          string `json:"nickname" validate:"required"`
	Email             string `json:"email" validate:"required"`
	Password          string `json:"password" validate:"required"`
	Language          string `json:"language" validate:"oneof=zh-CN en-US ja-JP ko-KR"`
	IsSubscribed      bool   `json:"isSubscribed"` // is subscribed are optional
	VerificationCode  string `json:"verificationCode" validate:"required"`
	VerificationToken string `json:"verificationToken" validate:"required"`
	InviteToken       string `json:"inviteToken"` // invite token are optional
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
