package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type GetUserByIDResponse struct {
	ID               string    `json:"userID"`
	UID              uuid.UUID `json:"uid"`
	Nickname         string    `json:"nickname"`
	Email            string    `json:"email"`
	Avatar           string    `json:"avatar"`
	Language         string    `json:"language"`
	IsTutorialViewed bool      `json:"isTutorialViewed"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func NewGetUserByIDResponse(user *User) *GetUserByIDResponse {
	customization := user.ExportUserCustomization()
	resp := &GetUserByIDResponse{
		ID:               idconvertor.ConvertIntToString(user.ID),
		UID:              user.UID,
		Email:            user.Email,
		Avatar:           user.Avatar,
		Language:         customization.Language,
		IsTutorialViewed: customization.IsTutorialViewed,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}
	// fill nickname
	if user.Nickname == PENDING_USER_NICKNAME { // reset pending user name for frontend display
		resp.Nickname = ""
	} else {
		resp.Nickname = user.Nickname
	}
	return resp
}

func (resp *GetUserByIDResponse) ExportForFeedback() interface{} {
	return resp
}
