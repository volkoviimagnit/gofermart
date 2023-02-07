package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
)

func TestUserOrdersGETHandler_ServeHTTP_Positive(t *testing.T) {
	testEnvironment := NewTestEnvironment()
	needOrders := 5

	tests := []struct {
		name                   string
		sameUser               bool
		expectedError          error
		expectedStatusCode     int
		expectedReceivedOrders int
	}{
		{
			name:                   "проверка выдачи заказов - 200",
			sameUser:               true,
			expectedError:          nil,
			expectedStatusCode:     http.StatusOK,
			expectedReceivedOrders: needOrders,
		},
		{
			name:                   "нет данных для ответа - 204",
			sameUser:               false,
			expectedError:          io.EOF,
			expectedStatusCode:     http.StatusNoContent,
			expectedReceivedOrders: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessToken := testEnvironment.CreateAndAuthorizeRandomUser(t)
			createdOrderDTOs, errOrderCreating := testEnvironment.CreateUserOrders(accessToken, needOrders)
			assert.Equal(t, len(createdOrderDTOs), needOrders)
			if errOrderCreating != nil {
				return
			}

			if !tt.sameUser {
				accessToken = testEnvironment.CreateAndAuthorizeRandomUser(t)
			}

			jsonResponse := testEnvironment.ServeHandler(testEnvironment.userOrderGETHandler, []byte(""), accessToken)

			receivedOrderDTOs := make([]response.UserOrderDTO, 0)
			jsonDecoder := json.NewDecoder(jsonResponse.Body)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					assert.NoError(t, err)
				}
			}(jsonResponse.Body)
			errDecoding := jsonDecoder.Decode(&receivedOrderDTOs)
			assert.Equal(t, tt.expectedError, errDecoding)
			assert.Equal(t, tt.expectedReceivedOrders, len(receivedOrderDTOs))
			assert.Equal(t, tt.expectedStatusCode, jsonResponse.StatusCode)
		})
	}
}

func TestUserOrdersGETHandler_ServeHTTP_Negative(t *testing.T) {
	testEnvironment := NewTestEnvironment()
	needOrders := 1
	accessToken := testEnvironment.CreateAndAuthorizeRandomUser(t)
	createdOrderDTOs, errOrderCreating := testEnvironment.CreateUserOrders(accessToken, needOrders)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "пользователь не авторизован - 401 (пустой токен)",
			token: "",
		},
		{
			name:  "пользователь не авторизован - 401 (неверный токен)",
			token: helpers.RandomString(10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, len(createdOrderDTOs), needOrders)
			assert.NoError(t, errOrderCreating)
			jsonResponse := testEnvironment.ServeHandler(testEnvironment.userOrderGETHandler, []byte(""), tt.token)
			receivedOrderDTOs := make([]response.UserOrderDTO, 0)
			jsonDecoder := json.NewDecoder(jsonResponse.Body)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					assert.NoError(t, err)
				}
			}(jsonResponse.Body)
			errDecoding := jsonDecoder.Decode(&receivedOrderDTOs)
			assert.Equal(t, io.EOF, errDecoding)
			assert.Equal(t, 0, len(receivedOrderDTOs))
			assert.Equal(t, http.StatusUnauthorized, jsonResponse.StatusCode)
		})
	}

}
