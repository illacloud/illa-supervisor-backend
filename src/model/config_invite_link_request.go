package model

type ConfigInviteLinkRequest struct {
	InviteLinkEnabled bool `json:"inviteLinkEnabled"`
}

func NewConfigInviteLinkRequest() *ConfigInviteLinkRequest {
	return &ConfigInviteLinkRequest{}
}
