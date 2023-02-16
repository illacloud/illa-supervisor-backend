package model

type CreateUserRequest struct {
	Nickname       string `json:"nickname" gorm:"column:nickname;type:varchar;size:15;not null"`
	PasswordDigest string `json:"passworddigest" gorm:"column:password_digest;type:varchar;size:60;not null"`
	Email          string `json:"email" gorm:"column:email;type:varchar;size:255;not null"`
	SSOConfig      string `json:"SSOConfig" gorm:"column:sso_config;type:jsonb"`        // for single sign-on data
	Customization  string `json:"customization" gorm:"column:customization;type:jsonb"` // for user itself customization config, including: Language, IsSubscribed

}

func NewCreateUserRequest() *CreateUserRequest {
	return &CreateUserRequest{}
}
