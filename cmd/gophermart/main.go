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

	userRegisterHandler := handlers.NewUserRegisterHandler(userRepository)
	userLoginHandler := handlers.NewUserLoginHandler(userRepository)
	userOrderPOSTHandler := handlers.NewUserOrderPOSTHandler()
	userOrderGETHandler := handlers.NewUserOrdersGETHandler()
	userBalanceHandler := handlers.NewUserBalanceHandler()
	userBalanceWithdrawHandler := handlers.NewUserBalanceWithdrawHandler()
	userWithdrawalsHandler := handlers.NewUserWithdrawalsHandler()
	orderNumberHandler := handlers.NewOrderNumberHandler()

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
	router.Method(orderNumberHandler.GetMethod(), orderNumberHandler.GetPattern(), orderNumberHandler)

	listenShutDown()

	logrus.Fatal(http.ListenAndServe(params.GetRunAddress(), router))
}

func listenShutDown() {
	c := make(chan os.Signal, 0)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c // Blocks here until interrupted
		os.Exit(1)
	}()
}
