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
	*AbstractHandler
	ubwRepository repository.IUserBalanceWithdrawRepository
}

func NewUserWithdrawalsHandler(ubwRepository repository.IUserBalanceWithdrawRepository, auth security.IAuthenticator) *UserWithdrawalsHandler {
	abstract := NewAbstractHandler(http.MethodGet, "/api/user/withdrawals", "application/json")
	abstract.SetAuthenticator(auth)
	return &UserWithdrawalsHandler{
		AbstractHandler: abstract,
		ubwRepository:   ubwRepository,
	}
}

func (h *UserWithdrawalsHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusNoContent
	// http.StatusUnauthorized
	// http.StatusInternalServerError

	passport := h.AuthOrAbort(rw, request)
	if passport == nil {
		return
	}

	// TODO добавить сортировку по времени вывода от самых старых к самым новым
	userWithdrawals, errFinding := h.ubwRepository.FindByUserID(passport.GetUser().GetID())
	if errFinding != nil {
		h.RenderInternalServerError(rw, errFinding)
		return
	}
	if len(userWithdrawals) == 0 {
		h.RenderNoContent(rw)
		return
	}

	var tempDTO response.UserWithdrawalDTO
	DTOs := make([]response.UserWithdrawalDTO, 0)

	for _, userWithdrawal := range userWithdrawals {
		tempDTO = response.UserWithdrawalDTO{
			OrderNumber: userWithdrawal.OrderNumber,
			Sum:         userWithdrawal.Sum,
			ProcessedAt: userWithdrawal.ProcessedAt.Format(time.RFC3339),
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
