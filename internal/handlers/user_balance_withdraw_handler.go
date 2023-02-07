package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/security"
	"github.com/volkoviimagnit/gofermart/internal/service"
)

type UserBalanceWithdrawHandler struct {
	parent             *AbstractHandler
	userBalanceService service.IUserBalanceService
	auth               security.IAuthenticator
}

func NewUserBalanceWithdrawHandler(ubService service.IUserBalanceService, auth security.IAuthenticator) *UserBalanceWithdrawHandler {
	abstract := NewAbstractHandler(http.MethodPost, "/api/user/balance/withdraw", "application/json")
	abstract.SetAuthenticator(auth)
	return &UserBalanceWithdrawHandler{
		parent:             abstract,
		userBalanceService: ubService,
	}
}

func (h *UserBalanceWithdrawHandler) GetContentType() string {
	return h.parent.contentType
}

func (h *UserBalanceWithdrawHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserBalanceWithdrawHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserBalanceWithdrawHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// http.StatusOK
	// http.StatusUnauthorized
	// http.StatusPaymentRequired
	// http.StatusUnprocessableEntity
	// http.StatusInternalServerError

	passport := h.parent.AuthOrAbort(rw, r)
	if passport == nil {
		return
	}

	dto, errDTO := h.extractRequestDTO(r)
	if errDTO != nil {
		h.parent.RenderInternalServerError(rw, errDTO)
		return
	}

	errWithdrawing := h.userBalanceService.AddUserWithdraw(passport.GetUser().Id(), dto.GetOrderNumber(), dto.GetSum())
	if errWithdrawing == nil {
		h.parent.RenderResponse(rw, http.StatusOK, []byte("UserBalanceWithdrawHandler"))
		return
	}
	switch errWithdrawing.(type) {
	default:
		h.parent.RenderInternalServerError(rw, errWithdrawing)
		return
	case *service.NotEnoughFundsError:
		h.parent.RenderError(rw, http.StatusPaymentRequired, errWithdrawing)
		return
	case *service.IncorrectOrderNumberError:
		h.parent.RenderError(rw, http.StatusUnprocessableEntity, errWithdrawing)
		return
	}
}

func (h *UserBalanceWithdrawHandler) extractRequestDTO(r *http.Request) (*request.UserBalanceWithdrawDTO, error) {
	dto := request.NewUserBalanceWithdrawDTO()
	jsonDecoder := json.NewDecoder(r.Body)
	errDecoding := jsonDecoder.Decode(dto)
	if errDecoding != nil {
		return nil, errDecoding
	}
	return dto, nil
}
