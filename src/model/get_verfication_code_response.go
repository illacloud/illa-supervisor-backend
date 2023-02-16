package model

type GetVerificationCodeResponse struct {
	VerificationToken string `json:"verificationToken"`
}

func NewGetVerificationCodeResponse(token string) *GetVerificationCodeResponse {
	return &GetVerificationCodeResponse{
		VerificationToken: token,
	}
}

func (resp *GetVerificationCodeResponse) ExportForFeedback() interface{} {
	return resp
}
