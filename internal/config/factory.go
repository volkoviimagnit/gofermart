package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

func GetConfig() (IConfig, error) {
	envVars, envErr := LoadEnv()
	if envErr != nil {
		return nil, envErr
	}
	argVars, argErr := loadArgs()
	if argErr != nil {
		return nil, argErr
	}

	var runAddress string
	var databaseURI string
	var accrualSystemAddress string
	var logLevel string
	if envVars.IsRunAddressSet() {
		runAddress = envVars.RunAddress
	} else {
		runAddress = argVars.RunAddress
	}
	if envVars.IsDatabaseURISet() {
		databaseURI = envVars.DatabaseURI
	} else {
		databaseURI = argVars.DatabaseURI
	}
	if envVars.IsAccrualSystemAddressSet() {
		accrualSystemAddress = envVars.AccrualSystemAddress
	} else {
		accrualSystemAddress = argVars.AccrualSystemAddress
	}
	if envVars.IsLogLevelSet() {
		logLevel = envVars.LogLevel
	} else {
		logLevel = argVars.LogLevel
	}

	return NewConfig(runAddress, databaseURI, accrualSystemAddress, logLevel), nil
}

func LoadEnv() (*Env, error) {
	var fileNames []string
	if _, err := os.Stat(".env.local"); err == nil {
		fileNames = append(fileNames, ".env.local")
	}
	if _, err := os.Stat(".env"); err == nil {
		fileNames = append(fileNames, ".env")
	}
	if len(fileNames) > 0 {
		errEnvFile := godotenv.Load(fileNames...)
		if errEnvFile != nil {
			fmt.Println("Cant load env-file")
			panic(errEnvFile)
		}
	}

	cfg := NewEnv()
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadArgs() (*Env, error) {
	args := NewEnv()

	flagSet := flag.NewFlagSet("default", flag.ExitOnError)
	flagSet.StringVar(&args.RunAddress, "a", "", "адрес и порт запуска сервиса")
	flagSet.StringVar(&args.DatabaseURI, "d", "", "адрес подключения к базе данных")
	flagSet.StringVar(&args.AccrualSystemAddress, "r", "", "адрес системы расчёта начислений")
	flagSet.StringVar(&args.LogLevel, "ll", "", "уровень логирования")

	flagError := flagSet.Parse(os.Args[1:])
	if flagError != nil {
		return nil, flagError
	}

	return args, nil
}
