package model

type TeamIconUploadAddressResponse struct {
	UploadAddress string `json:"uploadAddress"`
}

func NewTeamIconUploadAddressResponse(presignedURL string) *TeamIconUploadAddressResponse {
	return &TeamIconUploadAddressResponse{
		UploadAddress: presignedURL,
	}
}

func (resp *TeamIconUploadAddressResponse) ExportForFeedback() interface{} {
	return resp
}
