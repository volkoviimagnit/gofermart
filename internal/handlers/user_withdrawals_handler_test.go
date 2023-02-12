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

func TestUserWithdrawalsHandler_ServeHTTP_Negative_Positive(t *testing.T) {
	testEnvironment := NewTestEnvironment()
	needOrders := 5

	tests := []struct {
		name                   string
		sameUser               bool
		expectedError          error
		expectedStatusCode     int
		expectedUserWithdrawal int
	}{
		{
			name:                   "проверка выдачи заказов - 200",
			sameUser:               true,
			expectedError:          nil,
			expectedStatusCode:     http.StatusOK,
			expectedUserWithdrawal: needOrders,
		},
		{
			name:                   "нет данных для ответа - 204",
			sameUser:               false,
			expectedError:          io.EOF,
			expectedStatusCode:     http.StatusNoContent,
			expectedUserWithdrawal: 0,
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

			for _, createdOrderDTO := range createdOrderDTOs {
				_, errWithdrawing := testEnvironment.CreateUserBalanceWithdraw(accessToken, createdOrderDTO.GetNumber(), 100)
				assert.NoError(t, errWithdrawing)
			}

			if !tt.sameUser {
				accessToken = testEnvironment.CreateAndAuthorizeRandomUser(t)
			}

			jsonResponse := testEnvironment.ServeHandler(testEnvironment.userWithdrawalsHandler, []byte(""), accessToken)
			err := jsonResponse.Body.Close()
			assert.NoError(t, err)

			receivedUserWithdrawalDTOs := make([]response.UserWithdrawalDTO, 0)
			jsonDecoder := json.NewDecoder(jsonResponse.Body)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					assert.NoError(t, err)
				}
			}(jsonResponse.Body)
			errDecoding := jsonDecoder.Decode(&receivedUserWithdrawalDTOs)
			for _, receivedUserWithdrawalDTO := range receivedUserWithdrawalDTOs {
				assert.NotNil(t, receivedUserWithdrawalDTO.ProcessedAt)
				assert.NotEqual(t, "", receivedUserWithdrawalDTO.ProcessedAt)
			}
			assert.Equal(t, tt.expectedError, errDecoding)
			assert.Equal(t, tt.expectedUserWithdrawal, len(receivedUserWithdrawalDTOs))
			assert.Equal(t, tt.expectedStatusCode, jsonResponse.StatusCode)
		})
	}
}

func TestUserWithdrawalsHandler_ServeHTTP_Other(t *testing.T) {
	testEnvironment := NewTestEnvironment()

	tests := []struct {
		name        string
		accessToken string
	}{
		{
			name:        "пользователь не авторизован - 401 (пустой токен)",
			accessToken: "",
		},
		{
			name:        "пользователь не авторизован - 401 (неверный токен)",
			accessToken: helpers.RandomString(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonResponse := testEnvironment.ServeHandler(testEnvironment.userWithdrawalsHandler, []byte(""), tt.accessToken)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					assert.NoError(t, err)
				}
			}(jsonResponse.Body)
			assert.Equal(t, http.StatusUnauthorized, jsonResponse.StatusCode)
		})
	}
}
