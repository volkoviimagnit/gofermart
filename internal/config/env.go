package config

type Env struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel             string `env:"LOG_LEVEL"`
}

func NewEnv() *Env {
	return &Env{}
}

func (p *Env) IsRunAddressSet() bool {
	return len(p.RunAddress) > 0
}

func (p *Env) IsDatabaseURISet() bool {
	return len(p.DatabaseURI) > 0
}

func (p *Env) IsAccrualSystemAddressSet() bool {
	return len(p.AccrualSystemAddress) > 0
}

func (p *Env) IsLogLevelSet() bool {
	return len(p.LogLevel) > 0
}
