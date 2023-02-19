package handlers

import (
	"io"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/security"
	"github.com/volkoviimagnit/gofermart/internal/service"
)

type UserOrdersPOSTHandler struct {
	*AbstractHandler
	userOrderService service.IUserOrderService
}

func NewUserOrderPOSTHandler(userOrderService service.IUserOrderService, auth security.IAuthenticator) *UserOrdersPOSTHandler {
	abstract := NewAbstractHandler(http.MethodPost, "/api/user/orders", "text/plain")
	abstract.SetAuthenticator(auth)
	return &UserOrdersPOSTHandler{
		AbstractHandler:  abstract,
		userOrderService: userOrderService,
	}
}

func (h *UserOrdersPOSTHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// http.StatusOK
	//http.StatusAccepted
	//http.StatusBadRequest
	//http.StatusUnauthorized
	//http.StatusConflict
	//http.StatusUnprocessableEntity
	// http.StatusInternalServerError
	passport := h.AuthOrAbort(rw, r)
	if passport == nil {
		return
	}

	resp := response.NewResponse("text/plain")

	dto, errBody := h.extractRequestDTO(r)
	if errBody != nil {
		h.RenderInternalServerError(rw, errBody)
		return
	}

	var errStatusCode int
	errValidation := dto.Validate()
	if errValidation != nil {
		switch errValidation.(type) {
		case *request.NumberFormatError:
			errStatusCode = http.StatusUnprocessableEntity
		case *request.NumberError:
			errStatusCode = http.StatusBadRequest
		}
	}
	if errValidation != nil {
		h.RenderResponse(rw, errStatusCode, []byte(errValidation.Error()))
		return
	}

	errOrderAdding := h.userOrderService.AddOrder(passport.GetUser().GetID(), dto.GetNumber())
	switch errOrderAdding.(type) {
	case *service.RepositoryError:
		h.RenderInternalServerError(rw, errOrderAdding)
		return
	case *service.DuplicatedOwnOrderError:
		h.RenderResponse(rw, http.StatusOK, []byte(""))
		return
	case *service.DuplicatedSomebodyElseOrderError:
		h.RenderResponse(rw, http.StatusConflict, []byte(""))
		return
	}

	resp.SetStatus(http.StatusAccepted).SetBody([]byte(dto.GetNumber()))
	h.Render(rw, resp)
}

func (h *UserOrdersPOSTHandler) extractRequestDTO(r *http.Request) (*request.UserOrdersPOSTDTO, error) {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		return nil, errBody
	}
	requestDTO := request.NewUserOrdersPOSTDTO(string(body))
	return requestDTO, nil
}
