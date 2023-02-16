package model

import(
	"fmt"
)

type UpdateTeamPermissionRequest struct {
	AllowEditorInvite           bool `json:"allowEditorInvite"`
	AllowViewerInvite           bool `json:"allowViewerInvite"`
	AllowEditorManageTeamMember bool `json:"allowEditorManageTeamMember"`
	AllowViewerManageTeamMember bool `json:"allowViewerManageTeamMember"`
}

func NewUpdateTeamPermissionRequest(rawReq map[string]interface{}) *UpdateTeamPermissionRequest {
	fmt.Printf("rawReq: %v\n", rawReq)
	return &UpdateTeamPermissionRequest{}
}
