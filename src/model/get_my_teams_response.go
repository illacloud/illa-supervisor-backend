package model

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type GetMyTeamsResponse struct {
	MyTeams []*MyTeamResponse
}

type MyTeamResponse struct {
	ID                   string                `json:"id"`
	UID                  uuid.UUID             `json:"uid"`
	Name                 string                `json:"name"`
	Identifier           string                `json:"identifier"`
	Icon                 string                `json:"icon"`
	MyRole               int                   `json:"myRole"`
	TeamMemberID         string                `json:"teamMemberID"`
	TeamMemberPermission *TeamMemberPermission `json:"teamMemberPermission"`
	TeamPermission       *TeamPermission       `json:"permission"`
	JoinedAt             time.Time             `json:"-"`
}

func NewMyTeamResponse(team *Team, targetTeamMember *TeamMember) *MyTeamResponse {
	return &MyTeamResponse{
		ID:                   idconvertor.ConvertIntToString(team.ID),
		UID:                  team.UID,
		Name:                 team.Name,
		Identifier:           team.Identifier,
		Icon:                 team.Icon,
		MyRole:               targetTeamMember.UserRole,
		TeamMemberID:         idconvertor.ConvertIntToString(targetTeamMember.ID),
		TeamMemberPermission: targetTeamMember.ExportPermission(),
		TeamPermission:       team.ExportTeamPermission(),
		JoinedAt:             targetTeamMember.CreatedAt,
	}
}

func (resp *MyTeamResponse) ExportForFeedback() interface{} {
	return resp
}

func NewGetMyTeamsResponse(teams []*Team, teamMembersLT map[int]*TeamMemberForExport) *GetMyTeamsResponse {
	// build team  members lookup table
	ret := &GetMyTeamsResponse{}
	for _, team := range teams {
		targetTeamMember := teamMembersLT[team.ID]
		myTeam := &MyTeamResponse{
			ID:                   idconvertor.ConvertIntToString(team.ID),
			UID:                  team.UID,
			Name:                 team.Name,
			Identifier:           team.Identifier,
			Icon:                 team.Icon,
			MyRole:               targetTeamMember.UserRole,
			TeamMemberID:         idconvertor.ConvertIntToString(targetTeamMember.ID),
			TeamMemberPermission: targetTeamMember.Permission,
			TeamPermission:       team.ExportTeamPermission(),
			JoinedAt:             targetTeamMember.CreatedAt,
		}
		ret.MyTeams = append(ret.MyTeams, myTeam)
	}
	// sort
	sort.Slice(ret.MyTeams, func(i, j int) bool {
		return ret.MyTeams[i].JoinedAt.After(ret.MyTeams[j].JoinedAt)
	})
	return ret
}

func (resp *GetMyTeamsResponse) ExportForFeedback() interface{} {
	return resp.MyTeams
}
