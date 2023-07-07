package model

import (
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type GenerateInviteLinkResponse struct {
	TeamID     string `json:"teamID"`
	AppID      string `json:"appID"`
	UserRole   int    `json:"userRole"`
	InviteLink string `json:"inviteLink"`
}

func NewGenerateInviteLinkResponse() *GenerateInviteLinkResponse {
	return &GenerateInviteLinkResponse{}
}

func (u *GenerateInviteLinkResponse) SetTeamID(teamID int) {
	u.TeamID = idconvertor.ConvertIntToString(teamID)
}

func (u *GenerateInviteLinkResponse) SetAppID(appID int) {
	u.AppID = idconvertor.ConvertIntToString(appID)
}

func (u *GenerateInviteLinkResponse) SetUserRole(userRole int) {
	u.UserRole = userRole
}

func (u *GenerateInviteLinkResponse) SetInviteLink(inviteLink string) {
	u.InviteLink = inviteLink
}

func NewGenerateInviteLinkResponseByInvite(invite *Invite, redirectPage string) *GenerateInviteLinkResponse {
	i := NewGenerateInviteLinkResponse()
	i.TeamID = idconvertor.ConvertIntToString(invite.TeamID)
	i.AppID = idconvertor.ConvertIntToString(invite.AppID)
	i.UserRole = invite.UserRole
	if invite.IsShareAppInvite() {
		i.InviteLink = invite.ExportShareAppLink()
	} else {
		i.InviteLink = invite.ExportInviteLink()
	}
	// add redirectPage suffix to link
	i.InviteLink += GenerateRedirectPageParam(redirectPage)
	return i
}

func (resp *GenerateInviteLinkResponse) ExportForFeedback() interface{} {
	return resp
}
