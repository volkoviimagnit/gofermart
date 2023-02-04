package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/server"
)

func TestUserRegisterHandler_ServeHTTP(t *testing.T) {

	type Expected struct {
		Status int
	}

	type Request struct {
		dto         request.UserDTO
		contentType string
	}

	randomLogin := helpers.RandomString(10)
	randomPassword := helpers.RandomString(10)
	tests := []struct {
		name     string
		request  Request
		expected Expected
	}{
		{
			name: "Корректная регистрация - 200",
			request: Request{
				contentType: "application/json",
				dto: request.UserDTO{
					Login:    randomLogin,
					Password: randomPassword,
				},
			},
			expected: Expected{
				Status: http.StatusOK,
			},
		},
		{
			name: "Без логина - 400",
			request: Request{
				contentType: "application/json",
				dto: request.UserDTO{
					Password: helpers.RandomString(10),
				},
			},
			expected: Expected{
				Status: http.StatusBadRequest,
			},
		},
		{
			name: "Без пароля - 400",
			request: Request{
				contentType: "application/json",
				dto: request.UserDTO{
					Login: helpers.RandomString(10),
				},
			},
			expected: Expected{
				Status: http.StatusBadRequest,
			},
		},
		{
			name: "Без тела - 400",
			request: Request{
				contentType: "application/json",
				dto:         request.UserDTO{},
			},
			expected: Expected{
				Status: http.StatusBadRequest,
			},
		},
		{
			name: "Повторная регистрация - 409",
			request: Request{
				contentType: "application/json",
				dto: request.UserDTO{
					Login:    randomLogin,
					Password: randomPassword,
				},
			},
			expected: Expected{
				Status: http.StatusConflict,
			},
		},
	}

	userRepository := repository.NewUserRepositoryMem()
	userRegisterHandler := NewUserRegisterHandler(userRepository)
	handlerCollection := server.NewHandlerCollection()
	handlerCollection.AddHandler(userRegisterHandler)

	router := server.NewRouterChi(handlerCollection, true)
	ts := httptest.NewServer(router.GetHandler())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, errMarshaling := json.Marshal(tt.request.dto)
			require.NoError(t, errMarshaling)

			testRequest, errRequest := http.NewRequest(userRegisterHandler.GetMethod(), ts.URL+"/api/user/register", bytes.NewBuffer(body))
			testRequest.Header.Set("Content-Type", tt.request.contentType)
			require.NoError(t, errRequest)

			testWriter := httptest.NewRecorder()
			userRegisterHandler.ServeHTTP(testWriter, testRequest)

			assert.Equal(t, tt.expected.Status, testWriter.Result().StatusCode)
		})
	}
	defer ts.Close()
}
