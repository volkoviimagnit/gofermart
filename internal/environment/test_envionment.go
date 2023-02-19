package environment

import (
	"github.com/volkoviimagnit/gofermart/internal/repository"
)

type TestEnvironment struct {
	MainEnvironment
}

func (t *TestEnvironment) ConfigureDatabase() error {
	return nil
}

func (t *TestEnvironment) ConfigureRepositories() error {
	t.UserRepository = repository.NewUserRepositoryMem()
	t.UserOrderRepository = repository.NewUserOrderRepositoryMem()
	t.UserBalanceRepository = repository.NewUserBalanceRepositoryMem()
	t.UserBalanceWithdrawRepository = repository.NewUserBalanceWithdrawRepositoryMem()

	return nil
}

func NewTestEnvironment() *TestEnvironment {
	return &TestEnvironment{
		MainEnvironment: MainEnvironment{
			databaseURI:    "0.0.0.0",
			accrualAddress: "0.0.0.0",
			isDebugMode:    true,
		},
	}
}
