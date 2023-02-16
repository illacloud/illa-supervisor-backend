package model

type UserAvatarUploadAddressResponse struct {
	UploadAddress string `json:"uploadAddress"`
}

func NewUserAvatarUploadAddressResponse(presignedURL string) *UserAvatarUploadAddressResponse {
	return &UserAvatarUploadAddressResponse{
		UploadAddress: presignedURL,
	}
}

func (resp *UserAvatarUploadAddressResponse) ExportForFeedback() interface{} {
	return resp
}
