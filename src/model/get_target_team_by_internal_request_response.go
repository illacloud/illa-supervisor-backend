package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type GetTargetTeamByInternalRequestResponse struct {
	ID         string    `json:"id"`
	UID        uuid.UUID `json:"uid"`
	Name       string    `json:"name"`
	Identifier string    `json:"identifier"`
	Icon       string    `json:"icon"`
	Permission string    `json:"permission"` // for team permission config
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewGetTargetTeamByInternalRequestResponse(team *Team) *GetTargetTeamByInternalRequestResponse {
	resp := &GetTargetTeamByInternalRequestResponse{
		ID:         idconvertor.ConvertIntToString(team.ID),
		UID:        team.UID,
		Name:       team.Name,
		Identifier: team.Identifier,
		Icon:       team.Icon,
		Permission: team.Permission,
		CreatedAt:  team.CreatedAt,
		UpdatedAt:  team.UpdatedAt,
	}
	return resp
}

func (resp *GetTargetTeamByInternalRequestResponse) ExportForFeedback() interface{} {
	return resp
}
