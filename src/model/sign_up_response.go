package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type SignUpResponse struct {
	ID           string    `json:"id"`
	UID          uuid.UUID `json:"uid"`
	TeamMemberID string    `json:"teamMemberID"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	Avatar       string    `json:"avatar"`
	Language     string    `json:"language"`
	IsSubscribed bool      `json:"isSubscribed"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewSignUpResponse(u *User) *SignUpResponse {
	customization := u.ExportUserCustomization()
	resp := &SignUpResponse{
		ID:           idconvertor.ConvertIntToString(u.ID),
		UID:          u.UID,
		Nickname:     u.Nickname,
		Email:        u.Email,
		Avatar:       u.Avatar,
		Language:     customization.Language,
		IsSubscribed: customization.IsSubscribed,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
	return resp
}

func (resp *SignUpResponse) ExportForFeedback() interface{} {
	return resp
}
