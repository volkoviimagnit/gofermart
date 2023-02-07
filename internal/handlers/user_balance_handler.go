package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/security"
	"github.com/volkoviimagnit/gofermart/internal/service"
)

type UserBalanceHandler struct {
	parent    *AbstractHandler
	ubService service.IUserBalanceService
}

func NewUserBalanceHandler(ubService service.IUserBalanceService, auth security.IAuthenticator) *UserBalanceHandler {
	abstract := NewAbstractHandler(http.MethodGet, "/api/user/balance", "application/json")
	abstract.SetAuthenticator(auth)
	return &UserBalanceHandler{
		ubService: ubService,
		parent:    abstract,
	}
}

func (h *UserBalanceHandler) GetContentType() string {
	return h.parent.contentType
}

func (h *UserBalanceHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserBalanceHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserBalanceHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusUnauthorized
	// http.StatusInternalServerError
	passport := h.parent.AuthOrAbort(rw, request)
	if passport == nil {
		return
	}

	userBalance, err := h.ubService.GetUserBalance(passport.GetUser().Id())
	if err != nil {
		return
	}

	userBalanceDTO := response.NewUserBalanceDTO(userBalance.GetCurrent(), userBalance.GetWithdrawn())

	body, errMarshaling := json.Marshal(userBalanceDTO)
	if errMarshaling != nil {
		h.parent.RenderInternalServerError(rw, errMarshaling)
		return
	}

	h.parent.RenderResponse(rw, http.StatusOK, body)
}
