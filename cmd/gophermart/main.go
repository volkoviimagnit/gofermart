package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/volkoviimagnit/gofermart/internal/config"
	"github.com/volkoviimagnit/gofermart/internal/environment"
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

	env := environment.NewMainEnvironment(
		params.GetDatabaseURI(),
		params.GetAccrualSystemAddress(),
		params.IsDebugMode(),
	)
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

	listenShutDown()

	// TODO: добавить обработку занятого порта
	router, err := env.GetRouter()
	if err != nil {
		logrus.Fatalf("ошибка получение роутера - %s", err)
	}
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
