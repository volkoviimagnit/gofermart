package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/volkoviimagnit/gofermart/internal/config"
	"github.com/volkoviimagnit/gofermart/internal/handlers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/security"
	"github.com/volkoviimagnit/gofermart/internal/server"
	"github.com/volkoviimagnit/gofermart/internal/service"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	params, errConf := config.GetConfig()
	if errConf != nil {
		logrus.Fatalf("не удалось загрузить конфиг приложения: %+v", errConf)
	}
	logLevel, errLog := logrus.ParseLevel(params.GetLogLevel())
	if errLog != nil {
		logrus.Fatalf("не удалось получить уровень логгирования: %+v", errLog)
	}
	logrus.SetLevel(logLevel)

	logrus.Debugf("params: %+v", params)

	userRepository := repository.NewUserRepositoryMem()
	userOrderRepository := repository.NewUserOrderRepositoryMem()
	userBalanceRepository := repository.NewUserBalanceRepositoryMem()
	userBalanceWithdrawRepository := repository.NewUserBalanceWithdrawRepositoryMem()

	authenticator := security.NewAuthenticator(userRepository)

	userBalanceService := service.NewUserBalanceService(
		userBalanceRepository,
		userBalanceWithdrawRepository,
	)

	userRegisterHandler := handlers.NewUserRegisterHandler(userRepository)
	userLoginHandler := handlers.NewUserLoginHandler(userRepository, authenticator)
	userOrderPOSTHandler := handlers.NewUserOrderPOSTHandler(userOrderRepository, authenticator)
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

	listenShutDown()

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
