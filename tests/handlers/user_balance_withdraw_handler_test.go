package handlers

import (
	"math"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
)

func TestUserBalanceWithdrawHandler_ServeHTTP(t *testing.T) {

	testEnvironment := NewTestEnvironment()

	minimalUserBalanceWithdraw := request.NewUserBalanceWithdrawDTO()
	minimalUserBalanceWithdraw.OrderNumber = helpers.RandomOrderNumber()
	minimalUserBalanceWithdraw.Sum = 1

	maxUserBalanceWithdraw := request.NewUserBalanceWithdrawDTO()
	maxUserBalanceWithdraw.OrderNumber = helpers.RandomOrderNumber()
	maxUserBalanceWithdraw.Sum = math.MaxFloat64

	badUserBalanceWithdraw := request.NewUserBalanceWithdrawDTO()
	badUserBalanceWithdraw.OrderNumber = "1"
	badUserBalanceWithdraw.Sum = 1

	tests := []struct {
		name               string
		needAuth           bool
		dto                *request.UserBalanceWithdrawDTO
		expectedStatusCode int
	}{
		{
			name:               "успешная обработка запроса - 200",
			needAuth:           true,
			dto:                minimalUserBalanceWithdraw,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "пользователь не авторизован - 401",
			needAuth:           false,
			dto:                minimalUserBalanceWithdraw,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "на счету недостаточно средств - 402",
			needAuth:           true,
			dto:                maxUserBalanceWithdraw,
			expectedStatusCode: http.StatusPaymentRequired,
		},
		{
			name:               "неверный номер заказа - 422",
			needAuth:           true,
			dto:                badUserBalanceWithdraw,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessToken := ""
			if tt.needAuth {
				accessToken = testEnvironment.CreateAndAuthorizeRandomUser(t)
			}

			body, errSerializing := tt.dto.Serialize()
			require.NoError(t, errSerializing)

			response := testEnvironment.ServeHandler(testEnvironment.UserBalanceWithdrawHandler, body, accessToken)
			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)

			errClosing := response.Body.Close()
			assert.NoError(t, errClosing)
		})
	}
}
