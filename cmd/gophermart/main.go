package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/volkoviimagnit/gofermart/internal/config"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	params, errConf := config.GetConfig()
	if errConf != nil {
		logrus.Fatalf("не удалось загрузить конфиг приложения: %+v", errConf)
	}
	logLevel, errLog := logrus.ParseLevel(params.GetLogLevel())
	if errLog != nil {
		logrus.Fatalf("не получить уровень логгирования: %+v", errLog)
	}
	logrus.SetLevel(logLevel)

	logrus.Debugf("params: %+v", params)

	listenShutDown()
}

func listenShutDown() {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-termChan // Blocks here until interrupted
		os.Exit(1)
	}()
}
