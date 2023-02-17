package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type TeamMemberWithUserInfoForExportConverted struct {
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

func NewTeamMemberWithUserInfoForExportConverted(i *TeamMemberWithUserInfoForExport) *TeamMemberWithUserInfoForExportConverted {
	return &TeamMemberWithUserInfoForExportConverted{
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

type GetAllTeamMembersResponse struct {
	AllTeamMembers []*TeamMemberWithUserInfoForExportConverted
}

func NewGetAllTeamMembersResponse(d []*TeamMemberWithUserInfoForExport) *GetAllTeamMembersResponse {
	resp := &GetAllTeamMembersResponse{}
	for _, item := range d {
		resp.AllTeamMembers = append(resp.AllTeamMembers, NewTeamMemberWithUserInfoForExportConverted(item))
	}
	return resp
}

func (resp *GetAllTeamMembersResponse) ExportForFeedback() interface{} {
	return resp.AllTeamMembers
}
