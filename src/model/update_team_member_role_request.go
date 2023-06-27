package model

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
	if u.UserRole == USER_ROLE_OWNER {
		return true
	}
	return false
}
