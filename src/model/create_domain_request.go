package model

type CreateDomainRequest struct {
	TeamID             int    `json:"team_id" gorm:"column:team_id;type:bigserial;index:domains_team_id"`
	UserDomain         string `json:"user_domain" gorm:"column:user_domain;type:varchar;size:255;not null"`
	SystemDomainPrefix string `json:"system_domain_prefix" gorm:"column:system_domain_prefix;type:varchar;size:255;not null"`
	ResolveStatus      int    `json:"resolve_status" gorm:"column:resolve_status;type:smallint;"`
	Category           int    `json:"category" gorm:"column:category;type:smallint;"`
}

func NewCreateDomainRequest() *CreateDomainRequest {
	return &CreateDomainRequest{}
}
