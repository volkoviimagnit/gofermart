package environment

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/volkoviimagnit/gofermart/internal/client"
	"github.com/volkoviimagnit/gofermart/internal/db"
	"github.com/volkoviimagnit/gofermart/internal/handlers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/security"
	"github.com/volkoviimagnit/gofermart/internal/service"
	"github.com/volkoviimagnit/gofermart/internal/transport"
)

type MainEnvironment struct {
	databaseURI                   string
	accrualAddress                string
	isDebugMode                   bool
	dBConnection                  *db.ConnectionPostgres
	Authenticator                 *security.AuthenticatorHeader
	UserRepository                repository.IUserRepository
	UserOrderRepository           repository.IUserOrderRepository
	UserBalanceRepository         repository.IUserBalanceRepository
	UserBalanceWithdrawRepository repository.IUserBalanceWithdrawRepository
	messenger                     transport.IMessenger
	UserBalanceService            service.IUserBalanceService
	UserOrderService              service.IUserOrderService
	accrualHTTPClient             client.IAccrualClient
	UserRegisterHandler           *handlers.UserRegisterHandler
	UserLoginHandler              *handlers.UserLoginHandler
	UserOrderPOSTHandler          *handlers.UserOrdersPOSTHandler
	UserOrderGETHandler           *handlers.UserOrdersGETHandler
	UserBalanceHandler            *handlers.UserBalanceHandler
	UserBalanceWithdrawHandler    *handlers.UserBalanceWithdrawHandler
	UserWithdrawalsHandler        *handlers.UserWithdrawalsHandler
	router                        *chi.Mux
}

func (p *MainEnvironment) ConfigureDatabase() error {
	dbConnection := db.NewConnectionPostgres(context.Background(), p.databaseURI)
	dbConnectionError := dbConnection.TryConnect()
	if dbConnectionError != nil {
		return dbConnectionError
	}
	errMigrating := dbConnection.Migrate()
	if errMigrating != nil {
		return errMigrating
	}
	p.dBConnection = dbConnection
	return nil
}

func (p *MainEnvironment) ConfigureRepositories() error {
	if p.dBConnection == nil {
		return errors.New("невозможно сконфигурировать репозитории - нет соединения с БД")
	}

	p.UserRepository = repository.NewUserRepositoryPG(p.dBConnection)
	p.UserOrderRepository = repository.NewUserOrderRepositoryPG(p.dBConnection)
	p.UserBalanceRepository = repository.NewUserBalanceRepositoryPG(p.dBConnection)
	p.UserBalanceWithdrawRepository = repository.NewUserBalanceWithdrawRepositoryPG(p.dBConnection)
	return nil
}

func (p *MainEnvironment) ConfigureHTTPHandlers() error {
	repoError := p.checkRepositories()
	if repoError != nil {
		return fmt.Errorf("невозможно сконфигурировать обработчики - %s", repoError.Error())
	}

	serviceError := p.checkServices()
	if serviceError != nil {
		return fmt.Errorf("невозможно сконфигурировать обработчики - %s", serviceError.Error())
	}

	p.UserRegisterHandler = handlers.NewUserRegisterHandler(p.UserRepository, p.Authenticator)
	p.UserLoginHandler = handlers.NewUserLoginHandler(p.UserRepository, p.Authenticator)
	p.UserOrderPOSTHandler = handlers.NewUserOrderPOSTHandler(p.UserOrderService, p.Authenticator)
	p.UserOrderGETHandler = handlers.NewUserOrdersGETHandler(p.UserOrderRepository, p.Authenticator)
	p.UserBalanceHandler = handlers.NewUserBalanceHandler(p.UserBalanceService, p.Authenticator)
	p.UserBalanceWithdrawHandler = handlers.NewUserBalanceWithdrawHandler(p.UserBalanceService, p.Authenticator)
	p.UserWithdrawalsHandler = handlers.NewUserWithdrawalsHandler(p.UserBalanceWithdrawRepository, p.Authenticator)

	return nil
}

func (p *MainEnvironment) ConfigureRouter() error {
	handlerError := p.checkHandlers()
	if handlerError != nil {
		return fmt.Errorf("невозможно сконфигурировать роутер - %s", handlerError.Error())
	}

	p.router = chi.NewRouter()
	p.router.Use(middleware.RequestID)
	p.router.Use(middleware.RealIP)
	if p.isDebugMode {
		p.router.Use(middleware.Logger)
	}
	p.router.Use(middleware.Recoverer)
	p.router.Use(middleware.StripSlashes)
	p.router.Use(middleware.Timeout(60 * time.Second))

	p.router.Method(p.UserRegisterHandler.GetMethod(), p.UserRegisterHandler.GetPattern(), p.UserRegisterHandler)
	p.router.Method(p.UserLoginHandler.GetMethod(), p.UserLoginHandler.GetPattern(), p.UserLoginHandler)
	p.router.Method(p.UserOrderPOSTHandler.GetMethod(), p.UserOrderPOSTHandler.GetPattern(), p.UserOrderPOSTHandler)
	p.router.Method(p.UserOrderGETHandler.GetMethod(), p.UserOrderGETHandler.GetPattern(), p.UserOrderGETHandler)
	p.router.Method(p.UserBalanceHandler.GetMethod(), p.UserBalanceHandler.GetPattern(), p.UserBalanceHandler)
	p.router.Method(p.UserBalanceWithdrawHandler.GetMethod(), p.UserBalanceWithdrawHandler.GetPattern(), p.UserBalanceWithdrawHandler)
	p.router.Method(p.UserWithdrawalsHandler.GetMethod(), p.UserWithdrawalsHandler.GetPattern(), p.UserWithdrawalsHandler)

	return nil
}

func (p *MainEnvironment) ConfigureServices() error {
	repoError := p.checkRepositories()
	if repoError != nil {
		return errors.New("невозможно сконфигурировать сервисы - ошибка репозиториев")
	}

	p.Authenticator = security.NewAuthenticator(p.UserRepository)
	p.messenger = transport.NewMessengerMem()

	p.UserBalanceService = service.NewUserBalanceService(
		p.UserBalanceRepository,
		p.UserBalanceWithdrawRepository,
		p.UserOrderRepository,
		p.messenger,
	)

	if len(p.accrualAddress) == 0 {
		return errors.New("невозможно сконфигурировать сервисы - не задан адрес системы лояльности")
	}

	p.accrualHTTPClient = client.NewAccrualHTTPClient(p.accrualAddress)
	p.UserOrderService = service.NewUserOrderService(
		p.accrualHTTPClient,
		p.messenger,
		p.UserOrderRepository,
		p.UserBalanceRepository,
		p.UserBalanceWithdrawRepository,
	)

	return nil
}

// RunConsumers
// RunConsumers TODO кол-во конзюмеров вынести в переменные
func (p *MainEnvironment) RunConsumers() error {
	serviceError := p.checkServices()
	if serviceError != nil {
		return fmt.Errorf("невозможно сконфигурировать обработчики - %s", serviceError.Error())
	}
	for i := 0; i < 10; i++ {
		p.messenger.AddConsumer(service.NewOrderAccrualConsumer(p.messenger, p.accrualHTTPClient, p.UserOrderService))
	}
	for i := 0; i < 10; i++ {
		p.messenger.AddConsumer(service.NewUserBalanceRecalculateConsumer(p.UserBalanceService))
	}
	p.messenger.Consume()

	return nil
}

func (p *MainEnvironment) GetRouter() (http.Handler, error) {
	if p.router == nil {
		return nil, errors.New("роутер не сконфигурирован")
	}

	return p.router, nil
}

// checkRepositories TODO: разделить проверки
func (p *MainEnvironment) checkRepositories() error {
	if p.UserRepository == nil ||
		p.UserBalanceRepository == nil ||
		p.UserBalanceWithdrawRepository == nil ||
		p.UserOrderRepository == nil {
		return errors.New("один из репозиториев не сконфигурированы")
	}
	return nil
}

func (p *MainEnvironment) checkServices() error {
	if p.Authenticator == nil {
		return errors.New("аутентификатор не сконфигурирован")
	}
	if p.UserOrderService == nil {
		return errors.New("сервис заказов не сконфигурирован")
	}
	if p.UserBalanceService == nil {
		return errors.New("сервис баланса не сконфигурирован")
	}
	if p.messenger == nil {
		return errors.New("сервис транспорта не сконфигурирован")
	}
	if p.accrualHTTPClient == nil {
		return errors.New("клиент взаимодействия с системой лояльности не сконфигурирован")
	}

	return nil
}

func (p *MainEnvironment) checkHandlers() error {
	if p.UserRegisterHandler == nil {
		return errors.New("обработчик регистраций не задан")
	}
	if p.UserLoginHandler == nil {
		return errors.New("обработчик авторизации не задан")
	}
	if p.UserOrderPOSTHandler == nil {
		return errors.New("обработчик приема заказов не задан")
	}
	if p.UserOrderGETHandler == nil {
		return errors.New("обработчик списка заказов не задан")
	}
	if p.UserBalanceHandler == nil {
		return errors.New("обработчик баланса не задан")
	}
	if p.UserBalanceWithdrawHandler == nil {
		return errors.New("обработчик списания средств не задан")
	}
	if p.UserWithdrawalsHandler == nil {
		return errors.New("обработчик истории списаний не задан")
	}
	return nil
}

func NewMainEnvironment(databaseURI string, accrualAddress string, isDebugMode bool) IEnvironment {
	return &MainEnvironment{
		databaseURI:    databaseURI,
		accrualAddress: accrualAddress,
		isDebugMode:    isDebugMode,
	}
}
