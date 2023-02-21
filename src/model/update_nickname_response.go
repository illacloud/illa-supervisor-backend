package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type UpdateUserResponse struct {
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

func NewUpdateUserResponse(user *User) *UpdateUserResponse {
	customization := user.ExportUserCustomization()
	resp := &UpdateUserResponse{
		ID:           idconvertor.ConvertIntToString(user.ID),
		UID:          user.UID,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Language:     customization.Language,
		IsSubscribed: customization.IsSubscribed,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
	// fill nickname
	if user.Nickname == PENDING_USER_NICKNAME { // reset pending user name for frontend display
		resp.Nickname = ""
	} else {
		resp.Nickname = user.Nickname
	}
	return resp
}

func (resp *UpdateUserResponse) ExportForFeedback() interface{} {
	return resp
}
