package model

import "github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"

type InviteMemberByEmailResponse struct {
	Email        string `json:"email"`
	EmailStatus  bool   `json:"emailStatus"`
	UserRole     int    `json:"userRole"`
	TeamMemberID string `json:"teamMemberID"`
	AppID        string `json:"appID"`
	Feedback     string `json:"feedback"`
}

func NewInviteMemberByEmailResponse() *InviteMemberByEmailResponse {
	return &InviteMemberByEmailResponse{}
}

func (u *InviteMemberByEmailResponse) SetEmail(email string) {
	u.Email = email
}

func (u *InviteMemberByEmailResponse) SetUserRole(userRole int) {
	u.UserRole = userRole
}

func (u *InviteMemberByEmailResponse) SetEmailStatus(emailStatus bool) {
	u.EmailStatus = emailStatus
}

func (u *InviteMemberByEmailResponse) SetEmailStatusByInvite(i *Invite) {
	u.EmailStatus = i.ExportEmailStatus()
}

func (u *InviteMemberByEmailResponse) SetFeedback(f string) {
	u.Feedback = f
}

func NewInviteMemberByEmailResponseByRequest(req *InviteMemberByEmailRequest) *InviteMemberByEmailResponse {
	i := NewInviteMemberByEmailResponse()
	i.Email = req.Email
	i.UserRole = req.UserRole
	return i
}

func NewInviteMemberByEmailResponseByInviteRecord(i *Invite, feedback string) *InviteMemberByEmailResponse {
	resp := NewInviteMemberByEmailResponse()
	resp.Email = i.Email
	resp.UserRole = i.UserRole
	resp.TeamMemberID = idconvertor.ConvertIntToString(i.TeamMemberID)
	resp.AppID = idconvertor.ConvertIntToString(i.AppID)
	resp.EmailStatus = i.EmailStatus
	resp.Feedback = feedback
	return resp
}

func (resp *InviteMemberByEmailResponse) ExportForFeedback() interface{} {
	return resp
}
