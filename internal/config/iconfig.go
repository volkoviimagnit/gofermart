package config

type IConfig interface {
	GetRunAddress() string
	GetDatabaseURI() string
	GetAccrualSystemAddress() string
	GetLogLevel() string
}
