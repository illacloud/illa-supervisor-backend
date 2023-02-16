package model

type GetTeamPermissionResponse struct {
	AllowEditorInvite           bool `json:"allowEditorInvite"`
	AllowViewerInvite           bool `json:"allowViewerInvite"`
	AllowEditorManageTeamMember bool `json:"allowEditorManageTeamMember"`
	AllowViewerManageTeamMember bool `json:"allowViewerManageTeamMember"`
	InviteLinkEnabled           bool `json:"inviteLinkEnabled"`
}

func NewGetTeamPermissionResponse(team *Team) *GetTeamPermissionResponse {
	tp := team.ExportTeamPermission()
	resp := &GetTeamPermissionResponse{
		AllowEditorInvite:           tp.AllowEditorInvite,
		AllowViewerInvite:           tp.AllowViewerInvite,
		AllowEditorManageTeamMember: tp.AllowEditorManageTeamMember,
		AllowViewerManageTeamMember: tp.AllowViewerManageTeamMember,
		InviteLinkEnabled:           tp.InviteLinkEnabled,
	}
	return resp
}

func (resp *GetTeamPermissionResponse) ExportForFeedback() interface{} {
	return resp
}
