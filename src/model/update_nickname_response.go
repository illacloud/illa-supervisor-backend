package model

type UpdateNicknameResponse struct {
	Nickname string `json:"nickname"`
}

func NewUpdateNicknameResponseByUser(user *User) *UpdateNicknameResponse {
	return &UpdateNicknameResponse{
		Nickname: user.Nickname,
	}
}

func (resp *UpdateNicknameResponse) ExportForFeedback() interface{} {
	return resp
}
