package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserRegisterHandler struct {
	parent         *AbstractHandler
	userRepository repository.IUserRepository
}

func NewUserRegisterHandler(repository repository.IUserRepository) *UserRegisterHandler {
	return &UserRegisterHandler{
		parent:         NewAbstractHandler(http.MethodPost, "/api/user/register", "application/json"),
		userRepository: repository,
	}
}

func (h *UserRegisterHandler) GetContentType() string {
	return h.parent.contentType
}

func (h *UserRegisterHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserRegisterHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserRegisterHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// http.StatusOK
	// http.StatusConflict
	// http.StatusBadRequest
	// http.StatusInternalServerError

	// TODO: выдать токен авторизации сразу после регистрации Authorization
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

	user, errConflict := h.userRepository.FindOneByLogin(dto.GetLogin())
	if errConflict != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte("ошибка поиска по логину - " + errConflict.Error()))
		h.parent.Render(rw, resp)
		return
	}
	if user != nil {
		resp.SetStatus(http.StatusConflict).SetBody([]byte("логин уже занят"))
		h.parent.Render(rw, resp)
		return
	}

	user = &model.User{}
	user.SetLogin(dto.Login)
	user.SetPassword(dto.Password)
	errRepository := h.userRepository.Insert(*user)
	if errRepository != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(errRepository.Error()))
		h.parent.Render(rw, resp)
		return
	}

	resp.SetStatus(http.StatusOK).SetBody([]byte(""))
	h.parent.Render(rw, resp)
	return
}

func (h *UserRegisterHandler) extractRequestDTO(r *http.Request) (*request.UserDTO, error) {
	dto := &request.UserDTO{}
	jsonDecoder := json.NewDecoder(r.Body)
	errDecoding := jsonDecoder.Decode(dto)
	if errDecoding != nil {
		return nil, errDecoding
	}
	return dto, nil
}
