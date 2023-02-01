package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/repository"
)

type UserLoginHandler struct {
	parent         *AbstractHandler
	userRepository repository.IUserRepository
}

func NewUserLoginHandler(userRepository repository.IUserRepository) IHandler {
	return &UserLoginHandler{
		parent:         NewAbstractHandler(http.MethodPost, "/api/user/login"),
		userRepository: userRepository,
	}
}

func (h *UserLoginHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserLoginHandler) GetPattern() string {
	return h.parent.GetPattern()
}

// ServeHTTP todo научиться различать ошибки
func (h *UserLoginHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// http.StatusOK
	// http.StatusBadRequest
	// http.StatusUnauthorized
	// http.StatusInternalServerError

	resp := response.NewResponse("application/json")

	dto, errDTO := h.extractRequestDTO(r)
	if errDTO != nil {
		resp.SetStatus(http.StatusBadRequest).SetBody([]byte(errDTO.Error()))
		h.parent.Render(rw, resp)
		return
	}

	errValidation := dto.Validate()
	if errValidation != nil {
		resp.SetStatus(http.StatusBadRequest).SetBody([]byte(errValidation.Error()))
		h.parent.Render(rw, resp)
		return
	}

	user, errRepository := h.userRepository.GetOneByCredentials(dto.GetLogin(), dto.GetPassword())
	if errRepository != nil {
		resp.SetStatus(http.StatusUnauthorized).SetBody([]byte(errRepository.Error()))
		h.parent.Render(rw, resp)
		return
	}

	dto.Login = user.Id()
	body, errMarshaling := json.Marshal(dto)
	if errMarshaling != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(errMarshaling.Error()))
		h.parent.Render(rw, resp)
		return
	}

	resp.SetStatus(http.StatusOK).SetBody(body)
	h.parent.Render(rw, resp)
	return
}

func (h *UserLoginHandler) extractRequestDTO(r *http.Request) (*request.UserDTO, error) {
	dto := &request.UserDTO{}
	jsonDecoder := json.NewDecoder(r.Body)
	errDecoding := jsonDecoder.Decode(dto)
	if errDecoding != nil {
		return nil, errDecoding
	}
	return dto, nil
}
