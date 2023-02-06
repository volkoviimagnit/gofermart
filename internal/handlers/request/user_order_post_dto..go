package request

import (
	"strings"

	"github.com/ShiraazMoollatjie/goluhn"
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
		return &NumberError{}
	}
	errLuhn := goluhn.Validate(dto.GetNumber())
	if errLuhn != nil {
		return &NumberFormatError{errText: errLuhn.Error()}
	}
	// todo добавить проверку только на цифры
	return nil
}

func (dto *UserOrdersPOSTDTO) Serialize() ([]byte, error) {
	return []byte(dto.number), nil
}

type NumberFormatError struct {
	errText string
}

func (e *NumberFormatError) Error() string {
	return e.errText
}

type NumberError struct {
}

func (e *NumberError) Error() string {
	return "неверный формат запроса"
}
