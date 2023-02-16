package model

import (
	"github.com/illacloud/illa-supervisior-backend/src/utils/idconvertor"
)

type InviteMemberByEmailRequest struct {
	UserRole int    `json:"userRole" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	AppIDRaw string `json:"appID"` // optional
	AppID    int    `json:"-"`     // optional
}

func (u *InviteMemberByEmailRequest) Init() {
	// AppID can be empty when invite member join the team
	if len(u.AppIDRaw) > 0 {
		u.AppID = idconvertor.ConvertStringToInt(u.AppIDRaw)
	}
}

func (u *InviteMemberByEmailRequest) ExportUserRole() int {
	return u.UserRole
}

func (u *InviteMemberByEmailRequest) ExportEmail() string {
	return u.Email
}

func (u *InviteMemberByEmailRequest) ExportAppID() int {
	return u.AppID
}

func (u *InviteMemberByEmailRequest) IsShareAppInvite() bool {
	if u.AppID > 0 {
		return true
	}
	return false
}

func NewInviteMemberByEmailRequest() *InviteMemberByEmailRequest {
	return &InviteMemberByEmailRequest{}
}
