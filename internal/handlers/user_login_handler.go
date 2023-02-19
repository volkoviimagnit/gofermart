package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/security"
)

type UserLoginHandler struct {
	*AbstractHandler
	userRepository repository.IUserRepository
	auth           security.IAuthenticator
}

func NewUserLoginHandler(userRepository repository.IUserRepository, auth security.IAuthenticator) *UserLoginHandler {
	return &UserLoginHandler{
		AbstractHandler: NewAbstractHandler(http.MethodPost, "/api/user/login", "application/json"),
		userRepository:  userRepository,
		auth:            auth,
	}
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
		h.Render(rw, resp)
		return
	}

	errValidation := dto.Validate()
	if errValidation != nil {
		resp.SetStatus(http.StatusBadRequest).SetBody([]byte(errValidation.Error()))
		h.Render(rw, resp)
		return
	}

	user, errRepository := h.userRepository.FindOneByCredentials(dto.GetLogin(), dto.GetPassword())
	if errRepository != nil || user == nil {
		h.RenderUnauthorized(rw)
		return
	}

	accessToken := h.auth.CreateAuthenticatedToken()
	tokenDTO := response.NewUserLoginDTO(accessToken)
	user.Token = accessToken
	errUserUpdating := h.userRepository.Update(*user)
	if errUserUpdating != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(errUserUpdating.Error()))
		h.Render(rw, resp)
		return
	}

	body, errMarshaling := json.Marshal(tokenDTO)
	if errMarshaling != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(errMarshaling.Error()))
		h.Render(rw, resp)
		return
	}

	// TODO передать токен через заголовок Authorization
	resp.SetStatus(http.StatusOK).SetBody(body)
	h.auth.RenderAuthenticatedToken(rw, accessToken)
	h.Render(rw, resp)
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
