package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/security"
)

type UserWithdrawalsHandler struct {
	parent        *AbstractHandler
	ubwRepository repository.IUserBalanceWithdrawRepository
}

func NewUserWithdrawalsHandler(ubwRepository repository.IUserBalanceWithdrawRepository, auth security.IAuthenticator) *UserWithdrawalsHandler {
	abstract := NewAbstractHandler(http.MethodGet, "/api/user/withdrawals", "application/json")
	abstract.SetAuthenticator(auth)
	return &UserWithdrawalsHandler{
		parent:        abstract,
		ubwRepository: ubwRepository,
	}
}

func (h *UserWithdrawalsHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserWithdrawalsHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserWithdrawalsHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusNoContent
	// http.StatusUnauthorized
	// http.StatusInternalServerError

	passport := h.parent.AuthOrAbort(rw, request)
	if passport == nil {
		return
	}

	// TODO добавить сортировку по времени вывода от самых старых к самым новым
	userWithdrawals, errFinding := h.ubwRepository.FindByUserId(passport.GetUser().Id())
	if errFinding != nil {
		h.parent.RenderInternalServerError(rw, errFinding)
		return
	}
	if len(userWithdrawals) == 0 {
		h.parent.RenderNoContent(rw)
		return
	}

	var tempDTO response.UserWithdrawalDTO
	DTOs := make([]response.UserWithdrawalDTO, 0)

	for _, userWithdrawal := range userWithdrawals {
		tempDTO = response.UserWithdrawalDTO{
			OrderNumber: userWithdrawal.GetOrderNumber(),
			Sum:         userWithdrawal.GetSum(),
			ProcessedAt: userWithdrawal.GetProcessedAt().Format(time.RFC3339),
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
