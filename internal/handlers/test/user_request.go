package test

import "github.com/volkoviimagnit/gofermart/internal/handlers/request"

type UserRequest struct {
	DTO         request.UserDTO
	ContentType string
}
