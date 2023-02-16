package model

import (
	"github.com/google/uuid"
)

type UpdateInvitationCodeRequest struct {
	ID            int       `json:"id" gorm:"column:id;type:bigserial;primary_key;index:teams_ukey"`
	UID           uuid.UUID `json:"uid" gorm:"column:uid;type:uuid;not null;index:teams_ukey"`
	TeamID        int       `json:"teamID" gorm:"column:team_id;type:bigserial;index:invitation_code_team_id"`
	ConsumeStatus int       `json:"consumeStatus" gorm:"column:consume_status;type:smallint;"`
}

func NewUpdateInvitationCodeRequest() *UpdateInvitationCodeRequest {
	return &UpdateInvitationCodeRequest{}
}
