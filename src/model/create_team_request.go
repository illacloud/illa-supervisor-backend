package model

import "github.com/google/uuid"

type CreateTeamRequest struct {
	Name           string    `json:"name" validate:"required"`
	TeamDomain     string    `json:"teamDomain"`
	Identifier     string    `json:"identifier" validate:"required"`
	InvitationCode uuid.UUID `json:"invitationCode" validate:"required"`
}

func NewCreateTeamRequest() *CreateTeamRequest {
	return &CreateTeamRequest{}
}

func (req CreateTeamRequest) ExportInvitationCodeInUUID() uuid.UUID {
	return req.InvitationCode
}

func (req CreateTeamRequest) ExportTeamDomain() string {
	return req.TeamDomain
}

func (req CreateTeamRequest) ExportName() string {
	return req.Name
}

func (req CreateTeamRequest) ExportIdentifier() string {
	return req.Identifier
}
