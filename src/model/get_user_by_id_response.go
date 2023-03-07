package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type GetUserByIDResponse struct {
	ID        string    `json:"id"`
	UID       uuid.UUID `json:"uid"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewGetUserByIDResponse(user *User) *GetUserByIDResponse {
	customization := user.ExportUserCustomization()
	resp := &GetUserByIDResponse{
		ID:        idconvertor.ConvertIntToString(user.ID),
		UID:       user.UID,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Language:  customization.Language,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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
