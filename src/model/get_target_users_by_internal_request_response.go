package model

import "strconv"

type GetTargetUsersByInternalRequestResponse struct {
	Users map[string]*GetTargetUserByInternalRequestResponse `json:"users"`
}

func NewGetTargetUsersByInternalRequestResponse(users []*User) *GetTargetUsersByInternalRequestResponse {
	var getTargetUsersByInternalRequestResponse GetTargetUsersByInternalRequestResponse
	getTargetUsersByInternalRequestResponse.Users = make(map[string]*GetTargetUserByInternalRequestResponse, len(users))
	for _, user := range users {
		userResp := NewGetTargetUserByInternalRequestResponse(user)
		getTargetUsersByInternalRequestResponse.Users[strconv.Itoa(user.ID)] = userResp
	}
	return &getTargetUsersByInternalRequestResponse
}

func (resp *GetTargetUsersByInternalRequestResponse) ExportForFeedback() interface{} {
	return resp
}
