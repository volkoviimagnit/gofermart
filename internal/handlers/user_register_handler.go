package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
	"github.com/volkoviimagnit/gofermart/internal/security"
)

type UserRegisterHandler struct {
	*AbstractHandler
	userRepository repository.IUserRepository
}

func NewUserRegisterHandler(repository repository.IUserRepository, auth security.IAuthenticator) *UserRegisterHandler {
	abstract := NewAbstractHandler(http.MethodPost, "/api/user/register", "application/json")
	abstract.SetAuthenticator(auth)
	return &UserRegisterHandler{
		AbstractHandler: abstract,
		userRepository:  repository,
	}
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
		h.Render(rw, resp)
		return
	}

	errValidation := dto.Validate()
	if errValidation != nil {
		resp.SetStatus(http.StatusBadRequest).SetBody([]byte(errValidation.Error()))
		h.Render(rw, resp)
		return
	}

	user, errConflict := h.userRepository.FindOneByLogin(dto.GetLogin())
	if errConflict != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte("ошибка поиска по логину - " + errConflict.Error()))
		h.Render(rw, resp)
		return
	}
	if user != nil {
		resp.SetStatus(http.StatusConflict).SetBody([]byte("логин уже занят"))
		h.Render(rw, resp)
		return
	}

	accessToken := h.auth.CreateAuthenticatedToken()
	user = &model.User{
		Login:    dto.Login,
		Password: dto.Password,
		Token:    accessToken,
	}
	errRepository := h.userRepository.Insert(*user)
	if errRepository != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(errRepository.Error()))
		h.Render(rw, resp)
		return
	}

	resp.SetStatus(http.StatusOK).SetBody([]byte(""))
	h.auth.RenderAuthenticatedToken(rw, accessToken)
	h.Render(rw, resp)
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
