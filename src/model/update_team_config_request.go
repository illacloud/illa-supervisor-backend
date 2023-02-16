package model

import(
	"fmt"
)

type UpdateTeamConfigRequest struct {
	Name       string `json:"name" validate:"required"`
	Identifier string `json:"identifier" validate:"required"`
	Icon       string `json:"icon"`
}

func NewUpdateTeamConfigRequest(rawReq map[string]interface{}) *UpdateTeamConfigRequest {
	fmt.Printf("rawReq: %v\n", rawReq)
	return &UpdateTeamConfigRequest{}
}
