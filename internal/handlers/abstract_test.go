package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/test"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/security"
	"github.com/volkoviimagnit/gofermart/internal/server"
	"github.com/volkoviimagnit/gofermart/internal/service"
)

type UserTestCase struct {
	name     string
	request  test.UserRequest
	expected test.Expected
}

type TestEnvironment struct {
	authenticator        *security.AuthenticatorHeader
	userRegisterHandler  *UserRegisterHandler
	userLoginHandler     *UserLoginHandler
	userOrderPOSTHandler *UserOrdersPOSTHandler
	userOrderGETHandler  *UserOrdersGETHandler
	testServer           *httptest.Server
}

func (env *TestEnvironment) CreateAndAuthorizeRandomUser(t *testing.T) string {
	randomLogin, randomPassword := env.CreateRandomUser(t)

	body, errMarshaling := json.Marshal(request.UserDTO{
		Login:    randomLogin,
		Password: randomPassword,
	})
	require.NoError(t, errMarshaling)

	response := env.ServeHandler(env.userLoginHandler, body)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	errClosing := response.Body.Close()
	assert.NoError(t, errClosing)

	assert.NotEmpty(t, response.Header.Get("Authorization"))

	accessToken := response.Header.Get("Authorization")
	assert.NotEmpty(t, accessToken)
	return accessToken
}

func (env *TestEnvironment) CreateRandomUser(t *testing.T) (string, string) {
	randomLogin := helpers.RandomString(10)
	randomPassword := helpers.RandomString(10)

	body, errMarshaling := json.Marshal(request.UserDTO{
		Login:    randomLogin,
		Password: randomPassword,
	})
	require.NoError(t, errMarshaling)

	registerResponse := env.ServeHandler(env.userRegisterHandler, body)
	assert.Equal(t, http.StatusOK, registerResponse.StatusCode)

	errRegisterClosing := registerResponse.Body.Close()
	assert.NoError(t, errRegisterClosing)
	return randomLogin, randomPassword
}

func (env *TestEnvironment) CreateUserOrders(accessToken string, count int) ([]*request.UserOrdersPOSTDTO, error) {
	DTOs := make([]*request.UserOrdersPOSTDTO, 0, count)
	for i := 0; i < count; i++ {
		DTO := request.NewUserOrdersPOSTDTO(helpers.RandomOrderNumber())
		DTOs = append(DTOs, DTO)

		body, errSerializing := DTO.Serialize()
		if errSerializing != nil {
			return nil, errSerializing
		}

		response := env.ServeHandler(env.userOrderPOSTHandler, body, accessToken)
		if response.StatusCode != http.StatusAccepted {
			return nil, errors.New("не удалось создать заказ")
		}
	}
	return DTOs, nil
}

func (env *TestEnvironment) ServeHandler(handler server.IHttpHandler, body []byte, accessToken ...string) *http.Response {
	buffer := bytes.NewBuffer(body)
	testRequest, errRequest := http.NewRequest(handler.GetMethod(), env.testServer.URL+handler.GetPattern(), buffer)
	testRequest.Header.Set("Content-Type", handler.GetContentType())
	if len(accessToken) > 0 {
		token := accessToken[0]
		env.authenticator.SetAuthenticatedToken(testRequest, token)
	}
	// TODO добавить обработку ошибок
	if errRequest != nil {
		return nil
	}

	testWriter := httptest.NewRecorder()
	handler.ServeHTTP(testWriter, testRequest)

	httpResponse := testWriter.Result()

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(httpResponse.Body)

	return httpResponse
}

func NewTestEnvironment() *TestEnvironment {
	userRepository := repository.NewUserRepositoryMem()
	userOrderRepository := repository.NewUserOrderRepositoryMem()
	userBalanceRepository := repository.NewUserBalanceRepositoryMem()
	userBalanceWithdrawRepository := repository.NewUserBalanceWithdrawRepositoryMem()

	authenticator := security.NewAuthenticator(userRepository)

	userBalanceService := service.NewUserBalanceService(
		userBalanceRepository,
		userBalanceWithdrawRepository,
	)

	userRegisterHandler := NewUserRegisterHandler(userRepository)
	userLoginHandler := NewUserLoginHandler(userRepository, authenticator)
	userOrderPOSTHandler := NewUserOrderPOSTHandler(userOrderRepository, authenticator)
	userOrderGETHandler := NewUserOrdersGETHandler(userOrderRepository, authenticator)
	userBalanceHandler := NewUserBalanceHandler()
	userBalanceWithdrawHandler := NewUserBalanceWithdrawHandler(userBalanceService, authenticator)
	userWithdrawalsHandler := NewUserWithdrawalsHandler(userBalanceWithdrawRepository, authenticator)

	handlerCollection := server.NewHandlerCollection()
	handlerCollection.
		AddHandler(userRegisterHandler).
		AddHandler(userLoginHandler).
		AddHandler(userOrderPOSTHandler).
		AddHandler(userOrderGETHandler).
		AddHandler(userBalanceHandler).
		AddHandler(userBalanceWithdrawHandler).
		AddHandler(userWithdrawalsHandler)

	router := server.NewRouterChi(handlerCollection, true)
	ts := httptest.NewServer(router.GetHandler())

	return &TestEnvironment{
		authenticator:        authenticator,
		userRegisterHandler:  userRegisterHandler,
		userLoginHandler:     userLoginHandler,
		userOrderPOSTHandler: userOrderPOSTHandler,
		userOrderGETHandler:  userOrderGETHandler,
		testServer:           ts,
	}
}
