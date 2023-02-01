package response

type UserLoginDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func NewUserLoginDTO(accessToken string) *UserLoginDTO {
	return &UserLoginDTO{AccessToken: accessToken, TokenType: "Bearer"}
}
