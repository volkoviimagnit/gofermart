package request

import "errors"

type UserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (dto *UserDTO) Validate() error {
	if len(dto.Login) == 0 {
		return errors.New("поле login является обязательным")
	}
	if len(dto.Password) == 0 {
		return errors.New("поле password является обязательным")
	}
	return nil
}

func (dto *UserDTO) GetLogin() string {
	return dto.Login
}

func (dto *UserDTO) GetPassword() string {
	return dto.Password
}
