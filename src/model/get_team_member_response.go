package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type GetTeamMemberResponse struct {
	ID           string                `json:"userID"`
	UID          uuid.UUID             `json:"uid"`
	TeamMemberID string                `json:"teamMemberID"`
	Nickname     string                `json:"nickname"`
	Email        string                `json:"email"`
	Avatar       string                `json:"avatar"`
	Language     string                `json:"language"`
	IsSubscribed bool                  `json:"isSubscribed"`
	UserRole     int                   `json:"userRole"`
	Permission   *TeamMemberPermission `json:"permission"` // for user permission config
	UserStatus   int                   `json:"userStatus"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}

func NewGetTeamMemberResponse(i *TeamMemberWithUserInfoForExport) *GetTeamMemberResponse {
	return &GetTeamMemberResponse{
		ID:           idconvertor.ConvertIntToString(i.ID),
		UID:          i.UID,
		TeamMemberID: idconvertor.ConvertIntToString(i.TeamMemberID),
		Nickname:     i.Nickname,
		Email:        i.Email,
		Avatar:       i.Avatar,
		Language:     i.Language,
		IsSubscribed: i.IsSubscribed,
		UserRole:     i.UserRole,
		Permission:   i.Permission,
		UserStatus:   i.UserStatus,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
}

func (resp *GetTeamMemberResponse) ExportForFeedback() interface{} {
	return resp
}
