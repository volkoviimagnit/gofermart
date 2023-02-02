package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/security"
)

type UserOrdersGETHandler struct {
	parent       *AbstractHandler
	uoRepository repository.IUserOrderRepository
}

func NewUserOrdersGETHandler(uoRepository repository.IUserOrderRepository, auth security.IAuthenticator) *UserOrdersGETHandler {
	abstract := NewAbstractHandler(http.MethodGet, "/api/user/orders", "application/json")
	abstract.SetAuthenticator(auth)
	return &UserOrdersGETHandler{
		parent:       abstract,
		uoRepository: uoRepository,
	}
}

func (h *UserOrdersGETHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserOrdersGETHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserOrdersGETHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusNoContent
	// http.StatusUnauthorized
	// http.StatusInternalServerError

	passport := h.parent.AuthOrAbort(rw, request)
	if passport == nil {
		return
	}

	userOrders, errFinding := h.uoRepository.FindByUserId(passport.GetUser().Id())
	if errFinding != nil {
		h.parent.RenderInternalServerError(rw, errFinding)
		return
	}
	if len(userOrders) == 0 {
		h.parent.RenderNoContent(rw)
		return
	}

	var tempDTO response.UserOrderDTO
	DTOs := make([]response.UserOrderDTO, 0)

	for _, userOrder := range userOrders {
		tempDTO = response.UserOrderDTO{
			Number:     userOrder.Number(),
			Status:     userOrder.Status(),
			Accrual:    userOrder.Accrual(),
			UploadedAt: userOrder.UploadedAt().Format(time.RFC3339),
		}
		DTOs = append(DTOs, tempDTO)
	}

	body, errMarshaling := json.Marshal(DTOs)
	if errMarshaling != nil {
		h.parent.RenderInternalServerError(rw, errMarshaling)
		return
	}

	resp := response.NewResponse(h.parent.contentType)
	resp.SetStatus(http.StatusOK).SetBody(body)
	h.parent.Render(rw, resp)
	return
}
