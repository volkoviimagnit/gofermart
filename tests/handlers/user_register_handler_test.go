package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
	"github.com/volkoviimagnit/gofermart/tests/handlers/structs"
)

func TestUserRegisterHandler_ServeHTTP(t *testing.T) {
	randomLogin := helpers.RandomString(10)
	randomPassword := helpers.RandomString(10)
	tests := []UserTestCase{
		{
			name: "Корректная регистрация - 200",
			request: structs.UserRequest{
				DTO: request.UserDTO{
					Login:    randomLogin,
					Password: randomPassword,
				},
			},
			expected: structs.Expected{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Без логина - 400",
			request: structs.UserRequest{
				DTO: request.UserDTO{
					Password: helpers.RandomString(10),
				},
			},
			expected: structs.Expected{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Без пароля - 400",
			request: structs.UserRequest{
				DTO: request.UserDTO{
					Login: helpers.RandomString(10),
				},
			},
			expected: structs.Expected{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Без тела - 400",
			request: structs.UserRequest{
				DTO: request.UserDTO{},
			},
			expected: structs.Expected{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Повторная регистрация - 409",
			request: structs.UserRequest{
				DTO: request.UserDTO{
					Login:    randomLogin,
					Password: randomPassword,
				},
			},
			expected: structs.Expected{
				StatusCode: http.StatusConflict,
			},
		},
	}

	testEnvironment := NewTestEnvironment()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, errMarshaling := json.Marshal(tt.request.DTO)
			require.NoError(t, errMarshaling)

			response := testEnvironment.ServeHandler(testEnvironment.UserRegisterHandler, body)
			assert.Equal(t, tt.expected.StatusCode, response.StatusCode)

			errClosing := response.Body.Close()
			assert.NoError(t, errClosing)
		})
	}
}
