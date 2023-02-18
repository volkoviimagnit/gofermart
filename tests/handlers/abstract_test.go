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
	clientResponse "github.com/volkoviimagnit/gofermart/internal/client/response"
	"github.com/volkoviimagnit/gofermart/internal/environment"
	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
	"github.com/volkoviimagnit/gofermart/internal/server"
	"github.com/volkoviimagnit/gofermart/tests/handlers/structs"
)

type UserTestCase struct {
	name     string
	request  structs.UserRequest
	expected structs.Expected
}

type TestEnvironment struct {
	environment.TestEnvironment
	TestServer *httptest.Server
}

func (env *TestEnvironment) CreateAndAuthorizeRandomUser(t *testing.T) string {
	randomLogin, randomPassword := env.CreateRandomUser(t)

	body, errMarshaling := json.Marshal(request.UserDTO{
		Login:    randomLogin,
		Password: randomPassword,
	})
	require.NoError(t, errMarshaling)

	response := env.ServeHandler(env.UserLoginHandler, body)
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

	registerResponse := env.ServeHandler(env.UserRegisterHandler, body)
	assert.Equal(t, http.StatusOK, registerResponse.StatusCode)

	errRegisterClosing := registerResponse.Body.Close()
	assert.NoError(t, errRegisterClosing)

	user, errFindingUser := env.UserRepository.FindOneByLogin(randomLogin)
	assert.NoError(t, errFindingUser)
	userBalance, errBalancing := env.UserBalanceService.SetUserBalance(user.GetID(), math.MaxFloat32, 0)
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

		response := env.ServeHandler(env.UserOrderPOSTHandler, body, accessToken)
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
		err := env.UserOrderService.Update(DTO.GetNumber(), clientResponse.AccrualStatusProcessed, &accrual)
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

	response := env.ServeHandler(env.UserBalanceWithdrawHandler, body, accessToken)
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
	testRequest, errRequest := http.NewRequest(handler.GetMethod(), env.TestServer.URL+handler.GetPattern(), buffer)
	testRequest.Header.Set("Content-Type", handler.GetContentType())
	if len(accessToken) > 0 {
		token := accessToken[0]
		env.Authenticator.SetAuthenticatedToken(testRequest, token)
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
	env := environment.NewTestEnvironment()

	dbError := env.ConfigureDatabase()
	if dbError != nil {
		logrus.Fatalf("ошибка конфигурации БД - %s", dbError)
	}

	repoError := env.ConfigureRepositories()
	if repoError != nil {
		logrus.Fatalf("ошибка конфигурации репозиториев - %s", repoError)
	}

	serviceError := env.ConfigureServices()
	if serviceError != nil {
		logrus.Fatalf("ошибка конфигурации сервисов - %s", serviceError)
	}

	handlerError := env.ConfigureHTTPHandlers()
	if handlerError != nil {
		logrus.Fatalf("ошибка конфигурации обработчиков - %s", handlerError)
	}

	routerError := env.ConfigureRouter()
	if routerError != nil {
		logrus.Fatalf("ошибка конфигурации роутеров - %s", handlerError)
	}

	consumerError := env.RunConsumers()
	if consumerError != nil {
		logrus.Fatalf("ошибка запуска конзюмеров - %s", consumerError)
	}

	// TODO: добавить обработку занятого порта
	router, err := env.GetRouter()
	if err != nil {
		logrus.Fatalf("ошибка получение роутера - %s", err)
	}

	ts := httptest.NewServer(router)

	return &TestEnvironment{
		TestEnvironment: *env,
		TestServer:      ts,
	}
}
