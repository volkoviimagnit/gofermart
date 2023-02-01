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

func NewUserRegisterHandler(repository repository.IUserRepository) IHandler {
	return &UserRegisterHandler{
		parent:         NewAbstractHandler(http.MethodPost, "/api/user/register"),
		userRepository: repository,
	}
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

	user := model.User{}
	user.SetLogin(dto.Login)
	user.SetPassword(dto.Password)
	errRepository := h.userRepository.Insert(user)
	if errRepository != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(errRepository.Error()))
		h.parent.Render(rw, resp)
		return
	}

	entity, _ := h.userRepository.GetOneByCredentials(user.Login(), user.Password())
	dto.Login = entity.Id()

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

func (h *UserRegisterHandler) extractRequestDTO(r *http.Request) (*request.UserDTO, error) {
	dto := &request.UserDTO{}
	jsonDecoder := json.NewDecoder(r.Body)
	errDecoding := jsonDecoder.Decode(dto)
	if errDecoding != nil {
		return nil, errDecoding
	}
	return dto, nil
}
