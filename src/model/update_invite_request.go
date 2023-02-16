package model

import (
	"github.com/google/uuid"
)

type UpdateInviteRequest struct {
	ID       int       `json:"id" gorm:"column:id;type:bigserial;primary_key;index:invite_ukey"`
	UID      uuid.UUID `json:"uid" gorm:"column:uid;type:uuid;not null;index:invite_uid"`
	TeamID   int       `json:"teamID" gorm:"column:team_id;type:bigserial;index:users_ukey"`
	Email    string    `json:"email" gorm:"column:email;type:varchar;size:255;;index:invite_email"`
	UserRole int       `json:"userRole" gorm:"column:user_role;type:smallint;index:invite_user_role"`
	Status   int       `json:"status" gorm:"column:status;type:smallint;index:users_ukey"`
}

func NewUpdateInviteRequest() *UpdateInviteRequest {
	return &UpdateInviteRequest{}
}
