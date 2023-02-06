package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/test"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/server"
)

func TestUserRegisterHandler_ServeHTTP(t *testing.T) {
	randomLogin := helpers.RandomString(10)
	randomPassword := helpers.RandomString(10)
	tests := []UserTestCase{
		{
			name: "Корректная регистрация - 200",
			request: test.UserRequest{
				DTO: request.UserDTO{
					Login:    randomLogin,
					Password: randomPassword,
				},
			},
			expected: test.Expected{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Без логина - 400",
			request: test.UserRequest{
				DTO: request.UserDTO{
					Password: helpers.RandomString(10),
				},
			},
			expected: test.Expected{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Без пароля - 400",
			request: test.UserRequest{
				DTO: request.UserDTO{
					Login: helpers.RandomString(10),
				},
			},
			expected: test.Expected{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Без тела - 400",
			request: test.UserRequest{
				DTO: request.UserDTO{},
			},
			expected: test.Expected{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Повторная регистрация - 409",
			request: test.UserRequest{
				DTO: request.UserDTO{
					Login:    randomLogin,
					Password: randomPassword,
				},
			},
			expected: test.Expected{
				StatusCode: http.StatusConflict,
			},
		},
	}

	userRepository := repository.NewUserRepositoryMem()
	userRegisterHandler := NewUserRegisterHandler(userRepository)
	handlerCollection := server.NewHandlerCollection()
	handlerCollection.AddHandler(userRegisterHandler)

	router := server.NewRouterChi(handlerCollection, true)
	ts := httptest.NewServer(router.GetHandler())

	testEnvironment := NewTestEnvironment()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, errMarshaling := json.Marshal(tt.request.DTO)
			require.NoError(t, errMarshaling)
			registerResponse := testEnvironment.ServeHandler(testEnvironment.userRegisterHandler, body)
			assert.Equal(t, tt.expected.StatusCode, registerResponse.StatusCode)
			errClosing := registerResponse.Body.Close()
			assert.NoError(t, errClosing)
		})
	}
	defer ts.Close()
}
