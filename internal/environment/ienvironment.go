package environment

import (
	"net/http"
)

type IEnvironment interface {
	ConfigureDatabase() error
	ConfigureRepositories() error
	ConfigureServices() error
	ConfigureRouter() error
	ConfigureHTTPHandlers() error
	GetRouter() (http.Handler, error)
	RunConsumers() error

	//SetUserRepository(repo repository.IUserRepository)
	//SetUserRepository(repo repository.IUserRepository)
	//SetUserRepository(repo repository.IUserRepository)
	//SetUserRepository(repo repository.IUserRepository)
}
