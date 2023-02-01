package request

import (
	"errors"
	"strings"
)

type UserOrdersPOSTDTO struct {
	number string
}

func NewUserOrdersPOSTDTO(number string) *UserOrdersPOSTDTO {
	return &UserOrdersPOSTDTO{number: strings.TrimSpace(number)}
}

func (dto *UserOrdersPOSTDTO) GetNumber() string {
	return dto.number
}

func (dto *UserOrdersPOSTDTO) Validate() error {
	if len(dto.GetNumber()) == 0 {
		return errors.New("номер заказа является обязательным")
	}
	// todo добавить проверку только на цифры
	return nil
}
