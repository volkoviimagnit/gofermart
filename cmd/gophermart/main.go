package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/volkoviimagnit/gofermart/internal/config"
	"github.com/volkoviimagnit/gofermart/internal/handlers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/security"
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
	userBalanceHandler := handlers.NewUserBalanceHandler()
	userBalanceWithdrawHandler := handlers.NewUserBalanceWithdrawHandler(userBalanceService, authenticator)
	userWithdrawalsHandler := handlers.NewUserWithdrawalsHandler(userBalanceWithdrawRepository, authenticator)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	if params.IsDebugMode() {
		router.Use(middleware.Logger)
	}
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Method(userRegisterHandler.GetMethod(), userRegisterHandler.GetPattern(), userRegisterHandler)
	router.Method(userLoginHandler.GetMethod(), userLoginHandler.GetPattern(), userLoginHandler)
	router.Method(userOrderPOSTHandler.GetMethod(), userOrderPOSTHandler.GetPattern(), userOrderPOSTHandler)
	router.Method(userOrderGETHandler.GetMethod(), userOrderGETHandler.GetPattern(), userOrderGETHandler)
	router.Method(userBalanceHandler.GetMethod(), userBalanceHandler.GetPattern(), userBalanceHandler)
	router.Method(userBalanceWithdrawHandler.GetMethod(), userBalanceWithdrawHandler.GetPattern(), userBalanceWithdrawHandler)
	router.Method(userWithdrawalsHandler.GetMethod(), userWithdrawalsHandler.GetPattern(), userWithdrawalsHandler)

	listenShutDown()

	logrus.Fatal(http.ListenAndServe(params.GetRunAddress(), router))
}

func listenShutDown() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c // Blocks here until interrupted
		os.Exit(1)
	}()
}
