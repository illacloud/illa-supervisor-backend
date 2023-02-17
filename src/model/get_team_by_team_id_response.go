package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type GetTeamByTeamIDResponse struct {
	ID         string          `json:"id"`
	UID        uuid.UUID       `json:"uid"`
	Name       string          `json:"name"`
	Identifier string          `json:"identifier"`
	Icon       string          `json:"icon"`
	Permission *TeamPermission `json:"permission"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

func NewGetTeamByTeamIDResponse(t *Team) *GetTeamByTeamIDResponse {
	return &GetTeamByTeamIDResponse{
		ID:         idconvertor.ConvertIntToString(t.ID),
		UID:        t.UID,
		Name:       t.Name,
		Identifier: t.Identifier,
		Icon:       t.Icon,
		Permission: t.ExportTeamPermission(),
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func (resp *GetTeamByTeamIDResponse) ExportForFeedback() interface{} {
	return resp
}
