package model

import "github.com/illacloud/illa-supervisor-backend/src/accesscontrol"

type UpdateTeamMemberRoleRequest struct {
	UserRole int `json:"userRole" validate:"required"`
}

func NewUpdateTeamMemberRoleRequest() *UpdateTeamMemberRoleRequest {
	return &UpdateTeamMemberRoleRequest{}
}

func (u *UpdateTeamMemberRoleRequest) ExportUserRole() int {
	return u.UserRole
}

func (u *UpdateTeamMemberRoleRequest) IsTransferOwner() bool {
	if u.UserRole == accesscontrol.USER_ROLE_OWNER {
		return true
	}
	return false
}
