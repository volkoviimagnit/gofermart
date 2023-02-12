package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
)

func TestUserOrdersPOSTHandler_ServeHTTP_Negative(t *testing.T) {

	testEnvironment := NewTestEnvironment()

	tests := []struct {
		name               string
		dto                *request.UserOrdersPOSTDTO
		expectedStatusCode int
		needToken          bool
	}{
		{
			name:               "неверный формат запроса - пустой номер заказа- 400",
			dto:                request.NewUserOrdersPOSTDTO(""),
			expectedStatusCode: http.StatusBadRequest,
			needToken:          true,
		},
		{
			name:               "неверный формат запроса - буквы в номере заказа - 422",
			dto:                request.NewUserOrdersPOSTDTO(helpers.RandomString(10)),
			expectedStatusCode: http.StatusUnprocessableEntity,
			needToken:          true,
		},
		{
			name:               "неверный формат номера заказа - 422",
			dto:                request.NewUserOrdersPOSTDTO(helpers.RandomDigits(10)),
			expectedStatusCode: http.StatusUnprocessableEntity,
			needToken:          true,
		},
		{
			name:               "пользователь не аутентифицирован - 401",
			dto:                request.NewUserOrdersPOSTDTO(helpers.RandomString(10)),
			expectedStatusCode: http.StatusUnauthorized,
			needToken:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var accessToken string
			if tt.needToken {
				accessToken = testEnvironment.CreateAndAuthorizeRandomUser(t)
			}

			body, errSerializing := tt.dto.Serialize()
			require.NoError(t, errSerializing)

			response := testEnvironment.ServeHandler(testEnvironment.userOrderPOSTHandler, body, accessToken)
			err := response.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)
		})
	}
}

func TestUserOrdersPOSTHandler_ServeHTTP_DuplicatedNumber(t *testing.T) {
	testEnvironment := NewTestEnvironment()

	t.Run("номер заказа уже был загружен другим пользователем - 409", func(t *testing.T) {
		accessToken := testEnvironment.CreateAndAuthorizeRandomUser(t)
		orderDTOs, errOrderCreating := testEnvironment.CreateUserOrders(accessToken, 1)
		assert.NoError(t, errOrderCreating)

		anotherAccessToken := testEnvironment.CreateAndAuthorizeRandomUser(t)
		for _, orderDTO := range orderDTOs {
			body, errSerializing := orderDTO.Serialize()
			require.NoError(t, errSerializing)

			response := testEnvironment.ServeHandler(testEnvironment.userOrderPOSTHandler, body, anotherAccessToken)
			err := response.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, http.StatusConflict, response.StatusCode)
		}
	})
}

func TestUserOrdersPOSTHandler_ServeHTTP_Positive(t *testing.T) {
	testEnvironment := NewTestEnvironment()
	accessToken := testEnvironment.CreateAndAuthorizeRandomUser(t)
	orderNumber := helpers.RandomOrderNumber()
	tests := []struct {
		name               string
		expectedStatusCode int
	}{
		{name: "новый номер заказа принят в обработку - 202", expectedStatusCode: http.StatusAccepted},
		{name: "номер заказа уже был загружен этим пользователем - 200", expectedStatusCode: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderDTO := request.NewUserOrdersPOSTDTO(orderNumber)

			body, errMarshaling := orderDTO.Serialize()
			require.NoError(t, errMarshaling)

			response := testEnvironment.ServeHandler(testEnvironment.userOrderPOSTHandler, body, accessToken)
			err := response.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)
		})
	}
}
