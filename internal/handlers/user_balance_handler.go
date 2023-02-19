package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/security"
	"github.com/volkoviimagnit/gofermart/internal/service"
)

type UserBalanceHandler struct {
	*AbstractHandler
	ubService service.IUserBalanceService
}

func NewUserBalanceHandler(ubService service.IUserBalanceService, auth security.IAuthenticator) *UserBalanceHandler {
	abstract := NewAbstractHandler(http.MethodGet, "/api/user/balance", "application/json")
	abstract.SetAuthenticator(auth)
	return &UserBalanceHandler{
		ubService:       ubService,
		AbstractHandler: abstract,
	}
}

func (h *UserBalanceHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusUnauthorized
	// http.StatusInternalServerError
	passport := h.AuthOrAbort(rw, request)
	if passport == nil {
		return
	}

	userBalance, err := h.ubService.GetUserBalance(passport.GetUser().GetID())
	if err != nil {
		h.RenderInternalServerError(rw, err)
		return
	}

	userBalanceDTO := response.NewUserBalanceDTO(userBalance.GetCurrent(), userBalance.GetWithdrawn())

	body, errMarshaling := json.Marshal(userBalanceDTO)
	if errMarshaling != nil {
		h.RenderInternalServerError(rw, errMarshaling)
		return
	}

	h.RenderResponse(rw, http.StatusOK, body)
}
