package model

type UpdateTeamMemberRequest struct {
	ID         int    `json:"id" gorm:"column:id;type:bigserial;primary_key;index:team_members_ukey"`
	TeamID     int    `json:"teamID" gorm:"column:team_id;type:bigserial;index:team_members_team_and_user_id"`
	UserID     int    `json:"userID" gorm:"column:user_id;type:bigserial;index:team_members_team_and_user_id"`
	UserRole   int    `json:"userRole" gorm:"column:user_role;type:smallint"`
	Permission string `json:"permission" gorm:"column:authority;type:jsonb"` // for user permission config
}

func NewUpdateTeamMemberRequest() *UpdateTeamMemberRequest {
	return &UpdateTeamMemberRequest{}
}
