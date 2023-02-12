package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/volkoviimagnit/gofermart/internal/client"
	"github.com/volkoviimagnit/gofermart/internal/config"
	"github.com/volkoviimagnit/gofermart/internal/db"
	"github.com/volkoviimagnit/gofermart/internal/handlers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/security"
	"github.com/volkoviimagnit/gofermart/internal/server"
	"github.com/volkoviimagnit/gofermart/internal/service"
	"github.com/volkoviimagnit/gofermart/internal/transport"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	params, errConf := config.GetConfig(true)
	if errConf != nil {
		logrus.Fatalf("не удалось загрузить конфиг приложения: %+v", errConf)
	}
	logLevel, errLog := logrus.ParseLevel(params.GetLogLevel())
	if errLog != nil {
		logrus.Fatalf("не удалось получить уровень логгирования: %+v", errLog)
	}
	logrus.SetLevel(logLevel)

	logrus.Debugf("params: %+v", params)

	dbConnection := db.NewConnectionPostgres(context.Background(), params.GetDatabaseURI())
	dbConnectionError := dbConnection.TryConnect()
	if dbConnectionError != nil {
		logrus.Fatalf("ошибка соединения с БД - %s", dbConnectionError)
	}

	messenger := transport.NewMessengerMem()

	//userRepository := repository.NewUserRepositoryMem()
	userRepository := repository.NewUserRepositoryPG(dbConnection)
	userOrderRepository := repository.NewUserOrderRepositoryMem()
	//userBalanceRepository := repository.NewUserBalanceRepositoryMem()
	userBalanceRepository := repository.NewUserBalanceRepositoryPG(dbConnection)
	userBalanceWithdrawRepository := repository.NewUserBalanceWithdrawRepositoryMem()

	authenticator := security.NewAuthenticator(userRepository)

	userBalanceService := service.NewUserBalanceService(
		userBalanceRepository,
		userBalanceWithdrawRepository,
		userOrderRepository,
		messenger,
	)

	accrualHttpClient := client.NewAccrualHttpClient(params.GetAccrualSystemAddress())
	userOrderService := service.NewUserOrderService(
		accrualHttpClient,
		messenger,
		userOrderRepository,
		userBalanceRepository,
		userBalanceWithdrawRepository,
	)
	//userOrderService.AddOrder("1", "109")

	userRegisterHandler := handlers.NewUserRegisterHandler(userRepository)
	userLoginHandler := handlers.NewUserLoginHandler(userRepository, authenticator)
	userOrderPOSTHandler := handlers.NewUserOrderPOSTHandler(userOrderService, authenticator)
	userOrderGETHandler := handlers.NewUserOrdersGETHandler(userOrderRepository, authenticator)
	userBalanceHandler := handlers.NewUserBalanceHandler(userBalanceService, authenticator)
	userBalanceWithdrawHandler := handlers.NewUserBalanceWithdrawHandler(userBalanceService, authenticator)
	userWithdrawalsHandler := handlers.NewUserWithdrawalsHandler(userBalanceWithdrawRepository, authenticator)

	handlerCollection := server.NewHandlerCollection()
	handlerCollection.
		AddHandler(userRegisterHandler).
		AddHandler(userLoginHandler).
		AddHandler(userOrderPOSTHandler).
		AddHandler(userOrderGETHandler).
		AddHandler(userBalanceHandler).
		AddHandler(userBalanceWithdrawHandler).
		AddHandler(userWithdrawalsHandler)

	router := server.NewRouterChi(handlerCollection, params.IsDebugMode())
	err := router.Configure()
	if err != nil {
		logrus.Fatal("не удалось сконфигурировать роутер")
	}

	for i := 0; i < 10; i++ {
		messenger.AddConsumer(service.NewOrderAccrualConsumer(messenger, accrualHttpClient, userOrderService))
	}
	for i := 0; i < 10; i++ {
		messenger.AddConsumer(service.NewUserBalanceRecalculateConsumer(userBalanceService))
	}

	messenger.Consume(0, transport.OrderAccrualQueueName)

	listenShutDown()

	// TODO: добавить обработку занятого порта
	logrus.Fatal(http.ListenAndServe(params.GetRunAddress(), router.GetHandler()))
}

func listenShutDown() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c // Blocks here until interrupted
		os.Exit(1)
	}()
}
