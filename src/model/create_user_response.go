package model

import (
	"github.com/illacloud/illa-supervisor-backend/src/utils/idconvertor"
)

type CreateUserResponse struct {
	ID string `json:"id"`
}

func NewCreateUserResponse(id int) *CreateUserResponse {
	resp := &CreateUserResponse{
		ID: idconvertor.ConvertIntToString(id),
	}
	return resp
}

func (resp *CreateUserResponse) ExportForFeedback() interface{} {
	return resp
}
