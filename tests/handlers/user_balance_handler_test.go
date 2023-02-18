package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/tests/handlers/structs"
)

func TestUserBalanceHandler_ServeHTTP(t *testing.T) {
	testEnvironment := NewTestEnvironment()

	tests := []struct {
		name        string
		expected    structs.Expected
		needAuth    bool
		errDecoding error
	}{
		{
			name:        "успешная обработка запроса - 200",
			expected:    structs.Expected{StatusCode: http.StatusOK},
			needAuth:    true,
			errDecoding: nil,
		},
		{
			name:        "пользователь не авторизован - 401",
			expected:    structs.Expected{StatusCode: http.StatusUnauthorized},
			needAuth:    false,
			errDecoding: io.EOF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessToken := ""
			if tt.needAuth {
				accessToken = testEnvironment.CreateAndAuthorizeRandomUser(t)
			}

			jsonResponse := testEnvironment.ServeHandler(testEnvironment.UserBalanceHandler, []byte(""), accessToken)
			err := jsonResponse.Body.Close()
			assert.NoError(t, err)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					assert.NoError(t, err)
				}
			}(jsonResponse.Body)
			userBalanceDTO := response.UserBalanceDTO{}
			jsonDecoder := json.NewDecoder(jsonResponse.Body)
			errDecoding := jsonDecoder.Decode(&userBalanceDTO)
			assert.Equal(t, tt.errDecoding, errDecoding)
			assert.Equal(t, tt.expected.StatusCode, jsonResponse.StatusCode)
		})
	}
}
