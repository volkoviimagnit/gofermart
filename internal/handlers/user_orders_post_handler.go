package handlers

import (
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
	"github.com/volkoviimagnit/gofermart/internal/security"
)

type UserOrdersPOSTHandler struct {
	parent              *AbstractHandler
	userOrderRepository repository.IUserOrderRepository
}

func NewUserOrderPOSTHandler(uoRepository repository.IUserOrderRepository, auth security.IAuthenticator) *UserOrdersPOSTHandler {
	abstract := NewAbstractHandler(http.MethodPost, "/api/user/orders", "text/plain")
	abstract.SetAuthenticator(auth)
	return &UserOrdersPOSTHandler{
		parent:              abstract,
		userOrderRepository: uoRepository,
	}
}

func (h *UserOrdersPOSTHandler) GetContentType() string {
	return h.parent.contentType
}

func (h *UserOrdersPOSTHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserOrdersPOSTHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserOrdersPOSTHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	//http.StatusAccepted
	//http.StatusBadRequest
	//http.StatusUnauthorized
	//http.StatusConflict
	//http.StatusUnprocessableEntity
	// http.StatusInternalServerError
	passport := h.parent.AuthOrAbort(rw, request)
	if passport == nil {
		return
	}

	resp := response.NewResponse("text/plain")

	dto, errBody := h.extractRequestDTO(request)
	if errBody != nil {
		h.parent.RenderInternalServerError(rw, errBody)
		return
	}

	errValidation := dto.Validate()
	if errValidation != nil {
		resp.SetStatus(http.StatusBadRequest).SetBody([]byte(errValidation.Error()))
		h.parent.Render(rw, resp)
		return
	}

	m := model.UserOrder{}
	m.SetNumber(dto.GetNumber())
	m.SetUserId(passport.GetUser().Id())
	m.SetUploadedAt(time.Now())
	if rand.Int()%2 == 0 {
		accrual := rand.Float64()
		m.SetAccrual(&accrual)
	}
	errInserting := h.userOrderRepository.Insert(m)
	if errInserting != nil {
		h.parent.RenderInternalServerError(rw, errInserting)
		return
	}

	resp.SetStatus(http.StatusAccepted).SetBody([]byte(dto.GetNumber()))
	h.parent.Render(rw, resp)
	return
}

func (h *UserOrdersPOSTHandler) extractRequestDTO(r *http.Request) (*request.UserOrdersPOSTDTO, error) {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		return nil, errBody
	}
	requestDTO := request.NewUserOrdersPOSTDTO(string(body))
	return requestDTO, nil
}
