package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/test"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
)

func TestUserLoginHandler_ServeHTTP_Negative(t *testing.T) {
	testEnvironment := NewTestEnvironment()

	ts := testEnvironment.testServer
	defer ts.Close()

	tests := getNegativeTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, errMarshaling := json.Marshal(tt.request.DTO)
			require.NoError(t, errMarshaling)

			response := testEnvironment.ServeHandler(t, testEnvironment.userLoginHandler, body)
			assert.Equal(t, tt.expected.StatusCode, response.StatusCode)

			errClosing := response.Body.Close()
			assert.NoError(t, errClosing)
		})
	}
}

func TestUserLoginHandler_ServeHTTP_Positive(t *testing.T) {
	testEnvironment := NewTestEnvironment()
	ts := testEnvironment.testServer
	defer ts.Close()

	testEnvironment.CreateAndAuthorizeRandomUser(t)
}

func getNegativeTestCases() []UserTestCase {
	return []UserTestCase{
		{
			name: "неверная пара логин/пароль - 401",
			expected: test.Expected{
				StatusCode: http.StatusUnauthorized,
			},
			request: test.UserRequest{
				DTO: request.UserDTO{
					Login:    helpers.RandomString(10),
					Password: helpers.RandomString(10),
				},
				ContentType: "application/json",
			},
		},
		{
			name: "неверный формат запроса, без пароля - 400",
			expected: test.Expected{
				StatusCode: http.StatusBadRequest,
			},
			request: test.UserRequest{
				DTO: request.UserDTO{
					Login: "test",
				},
				ContentType: "application/json",
			},
		},
		{
			name: "неверный формат запроса, без логина - 400",
			expected: test.Expected{
				StatusCode: http.StatusBadRequest,
			},
			request: test.UserRequest{
				DTO: request.UserDTO{
					Password: "test",
				},
				ContentType: "application/json",
			},
		},
		{
			name: "неверный формат запроса, без логина/пароля - 400",
			expected: test.Expected{
				StatusCode: http.StatusBadRequest,
			},
			request: test.UserRequest{
				DTO:         request.UserDTO{},
				ContentType: "application/json",
			},
		},
	}
}
