package model

type CreateInviteRequest struct {
	TeamID       int    `json:"teamID" gorm:"column:team_id;type:bigserial;index:users_ukey"`
	Email        string `json:"email" gorm:"column:email;type:varchar;size:255;;index:invite_email"`
	UserRole     int    `json:"userRole" gorm:"column:user_role;type:smallint;index:invite_user_role"`
	InviteStatus int    `json:"inviteStatus" gorm:"column:invite_status;type:smallint;index:users_ukey"`
}

func NewCreateInviteRequest() *CreateInviteRequest {
	return &CreateInviteRequest{}
}
