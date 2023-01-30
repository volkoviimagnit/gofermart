package config

type Config struct {
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
	logLevel             string
}

func NewConfig(runAddress string, databaseURI string, accrualSystemAddress string, logLevel string) *Config {
	return &Config{runAddress: runAddress, databaseURI: databaseURI, accrualSystemAddress: accrualSystemAddress, logLevel: logLevel}
}

func (c *Config) GetLogLevel() string {
	return c.logLevel
}

func (c *Config) GetRunAddress() string {
	return c.runAddress
}

func (c *Config) GetDatabaseURI() string {
	return c.databaseURI
}

func (c *Config) GetAccrualSystemAddress() string {
	return c.accrualSystemAddress
}
