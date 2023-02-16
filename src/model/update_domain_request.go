package model

import (
	"github.com/google/uuid"
)

type UpdateDomainRequest struct {
	ID                 int       `json:"id" gorm:"column:id;type:bigserial;primary_key;index:domains_ukey"`
	UID                uuid.UUID `json:"uid" gorm:"column:uid;type:uuid;not null;index:domains_ukey"`
	TeamID             int       `json:"teamID" gorm:"column:team_id;type:bigserial;index:domains_team_id"`
	UserDomain         string    `json:"userDomain" gorm:"column:user_domain;type:varchar;size:255;not null"`
	SystemDomainPrefix string    `json:"systemDomain_prefix" gorm:"column:system_domain_prefix;type:varchar;size:255;not null"`
	ResolveStatus      int       `json:"resolveStatus" gorm:"column:resolve_status;type:smallint;"`
	Category           int       `json:"category" gorm:"column:category;type:smallint;"`
}

func NewUpdateDomainRequest() *UpdateDomainRequest {
	return &UpdateDomainRequest{}
}
