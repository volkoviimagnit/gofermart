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
	parent         *AbstractHandler
	userRepository repository.IUserRepository
	auth           security.IAuthenticator
}

func NewUserLoginHandler(userRepository repository.IUserRepository, auth security.IAuthenticator) *UserLoginHandler {
	return &UserLoginHandler{
		parent:         NewAbstractHandler(http.MethodPost, "/api/user/login", "application/json"),
		userRepository: userRepository,
		auth:           auth,
	}
}

func (h *UserLoginHandler) GetContentType() string {
	return h.parent.contentType
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

	user, errRepository := h.userRepository.FindOneByCredentials(dto.GetLogin(), dto.GetPassword())
	if errRepository != nil || user == nil {
		h.parent.RenderUnauthorized(rw)
		return
	}

	accessToken := h.auth.CreateAuthenticatedToken()
	tokenDTO := response.NewUserLoginDTO(accessToken)
	user.SetToken(accessToken)
	errUserUpdating := h.userRepository.Update(*user)
	if errUserUpdating != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(errUserUpdating.Error()))
		h.parent.Render(rw, resp)
		return
	}

	body, errMarshaling := json.Marshal(tokenDTO)
	if errMarshaling != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(errMarshaling.Error()))
		h.parent.Render(rw, resp)
		return
	}

	// TODO передать токен через заголовок Authorization
	resp.SetStatus(http.StatusOK).SetBody(body)
	h.auth.RenderAuthenticatedToken(rw, accessToken)
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
