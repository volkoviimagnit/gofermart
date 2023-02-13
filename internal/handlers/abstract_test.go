package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volkoviimagnit/gofermart/internal/client"
	clientResponse "github.com/volkoviimagnit/gofermart/internal/client/response"
	"github.com/volkoviimagnit/gofermart/internal/config"
	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/test"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/security"
	"github.com/volkoviimagnit/gofermart/internal/server"
	"github.com/volkoviimagnit/gofermart/internal/service"
	"github.com/volkoviimagnit/gofermart/internal/transport"
)

type UserTestCase struct {
	name     string
	request  test.UserRequest
	expected test.Expected
}

type TestEnvironment struct {
	authenticator              *security.AuthenticatorHeader
	userRegisterHandler        *UserRegisterHandler
	userLoginHandler           *UserLoginHandler
	userOrderPOSTHandler       *UserOrdersPOSTHandler
	userOrderGETHandler        *UserOrdersGETHandler
	userBalanceHandler         *UserBalanceHandler
	userBalanceWithdrawHandler *UserBalanceWithdrawHandler
	userWithdrawalsHandler     *UserWithdrawalsHandler
	userBalanceService         service.IUserBalanceService
	userOrderService           service.IUserOrderService
	userRepository             repository.IUserRepository
	testServer                 *httptest.Server
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

	user, errFindingUser := env.userRepository.FindOneByLogin(randomLogin)
	assert.NoError(t, errFindingUser)
	userBalance, errBalancing := env.userBalanceService.SetUserBalance(user.GetID(), math.MaxFloat32, 0)
	assert.NoError(t, errBalancing)
	assert.True(t, userBalance.GetCurrent() == math.MaxFloat32)

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
		errBodyClosing := response.Body.Close()
		if errBodyClosing != nil {
			return nil, errBodyClosing
		}
		if response.StatusCode != http.StatusAccepted {
			return nil, errors.New("не удалось создать заказ")
		}
	}

	for _, DTO := range DTOs {
		accrual := 10000.0
		err := env.userOrderService.Update(DTO.GetNumber(), clientResponse.AccrualStatusProcessed, &accrual)
		if err != nil {
			return nil, err
		}
	}

	return DTOs, nil
}

func (env *TestEnvironment) CreateUserBalanceWithdraw(accessToken string, orderNumber string, sum float64) (*request.UserBalanceWithdrawDTO, error) {
	userBalanceWithdrawDTO := request.NewUserBalanceWithdrawDTO()
	userBalanceWithdrawDTO.OrderNumber = orderNumber
	userBalanceWithdrawDTO.Sum = sum

	body, errSerializing := userBalanceWithdrawDTO.Serialize()
	if errSerializing != nil {
		return nil, errSerializing
	}

	response := env.ServeHandler(env.userBalanceWithdrawHandler, body, accessToken)
	errBodyClosing := response.Body.Close()
	if errBodyClosing != nil {
		return nil, errBodyClosing
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось создать списание средств - code: %s", response.Status)
	}
	return userBalanceWithdrawDTO, nil
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
	err := httpResponse.Body.Close()
	if err != nil {
		return nil
	}
	return httpResponse
}

func NewTestEnvironment() *TestEnvironment {
	params, errConf := config.GetConfig(false)
	if errConf != nil {
		logrus.Fatalf("ошибка cбора настроек - %s", errConf)
	}
	messenger := transport.NewMessengerMem()

	userRepository := repository.NewUserRepositoryMem()
	userOrderRepository := repository.NewUserOrderRepositoryMem()
	userBalanceRepository := repository.NewUserBalanceRepositoryMem()
	userBalanceWithdrawRepository := repository.NewUserBalanceWithdrawRepositoryMem()

	/*
		userRepository := repository.NewUserRepositoryPG(dbConnection)
		userOrderRepository := repository.NewUserOrderRepositoryPG(dbConnection)
		userBalanceRepository := repository.NewUserBalanceRepositoryPG(dbConnection)
		userBalanceWithdrawRepository := repository.NewUserBalanceWithdrawRepositoryPG(dbConnection)
	*/

	authenticator := security.NewAuthenticator(userRepository)

	userBalanceService := service.NewUserBalanceService(
		userBalanceRepository,
		userBalanceWithdrawRepository,
		userOrderRepository,
		messenger,
	)

	accrualHTTPClient := client.NewAccrualHTTPClient(params.GetAccrualSystemAddress())
	userOrderService := service.NewUserOrderService(
		accrualHTTPClient,
		messenger,
		userOrderRepository,
		userBalanceRepository,
		userBalanceWithdrawRepository,
	)
	//userOrderService.AddOrder("1", "109")

	userRegisterHandler := NewUserRegisterHandler(userRepository, authenticator)
	userLoginHandler := NewUserLoginHandler(userRepository, authenticator)
	userOrderPOSTHandler := NewUserOrderPOSTHandler(userOrderService, authenticator)
	userOrderGETHandler := NewUserOrdersGETHandler(userOrderRepository, authenticator)
	userBalanceHandler := NewUserBalanceHandler(userBalanceService, authenticator)
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
		authenticator:              authenticator,
		userRegisterHandler:        userRegisterHandler,
		userLoginHandler:           userLoginHandler,
		userOrderPOSTHandler:       userOrderPOSTHandler,
		userOrderGETHandler:        userOrderGETHandler,
		userBalanceHandler:         userBalanceHandler,
		userBalanceWithdrawHandler: userBalanceWithdrawHandler,
		userWithdrawalsHandler:     userWithdrawalsHandler,
		userBalanceService:         userBalanceService,
		userOrderService:           userOrderService,
		userRepository:             userRepository,
		testServer:                 ts,
	}
}
