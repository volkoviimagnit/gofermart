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
	*AbstractHandler
	uoRepository repository.IUserOrderRepository
}

func NewUserOrdersGETHandler(uoRepository repository.IUserOrderRepository, auth security.IAuthenticator) *UserOrdersGETHandler {
	abstract := NewAbstractHandler(http.MethodGet, "/api/user/orders", "application/json")
	abstract.SetAuthenticator(auth)
	return &UserOrdersGETHandler{
		AbstractHandler: abstract,
		uoRepository:    uoRepository,
	}
}

func (h *UserOrdersGETHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusNoContent
	// http.StatusUnauthorized
	// http.StatusInternalServerError

	passport := h.AuthOrAbort(rw, request)
	if passport == nil {
		return
	}

	userOrders, errFinding := h.uoRepository.FindByUserID(passport.GetUser().GetID())
	if errFinding != nil {
		h.RenderInternalServerError(rw, errFinding)
		return
	}
	if len(userOrders) == 0 {
		h.RenderNoContent(rw)
		return
	}

	var tempDTO response.UserOrderDTO
	DTOs := make([]response.UserOrderDTO, 0)

	for _, userOrder := range userOrders {
		tempDTO = response.UserOrderDTO{
			Number:     userOrder.Number,
			Status:     userOrder.Status.String(),
			Accrual:    userOrder.Accrual,
			UploadedAt: userOrder.UploadedAt.Format(time.RFC3339),
		}
		DTOs = append(DTOs, tempDTO)
	}

	body, errMarshaling := json.Marshal(DTOs)
	if errMarshaling != nil {
		h.RenderInternalServerError(rw, errMarshaling)
		return
	}

	resp := response.NewResponse(h.contentType)
	resp.SetStatus(http.StatusOK).SetBody(body)
	h.Render(rw, resp)
}
