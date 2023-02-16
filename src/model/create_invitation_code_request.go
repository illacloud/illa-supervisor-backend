package model

type CreateInvitationCodeRequest struct {
	TeamID        int `json:"teamID" gorm:"column:team_id;type:bigserial;index:invitation_code_team_id"`
	ConsumeStatus int `json:"consumeStatus" gorm:"column:consume_status;type:smallint;"`
}

func NewCreateInvitationCodeRequest() *CreateInvitationCodeRequest {
	return &CreateInvitationCodeRequest{}
}
