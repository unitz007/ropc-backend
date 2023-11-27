package utils

import (
	"log"
	"os"
	"strconv"
)

type Config interface {
	ServerPort() string
	DatabasePassword() string
	DatabaseUser() string
	DatabaseName() string
	TokenSecret() string
	DatabaseHost() string
	DatabasePort() string
	TokenExpiry() int
	NewRelicAppName() string
	NewRelicLicense() string
	Mux() string
	Environment() string
}

type config struct{}

func (e config) Environment() string {
	return getEnvironmentVariable("ENVIRONMENT")
}

func (e config) NewRelicAppName() string {
	return getEnvironmentVariable("NEW_RELIC_APP_NAME")
}

func (e config) NewRelicLicense() string {
	return getEnvironmentVariable("NEW_RELIC_LICENSE")
}

func (e config) Mux() string {
	return getEnvironmentVariable("ROPC_MUX")
}

func (e config) ServerPort() string {
	return getEnvironmentVariable("ROPC_SERVER_PORT")
}

func (e config) DatabasePassword() string {
	return getEnvironmentVariable("ROPC_DATABASE_PASSWORD")
}

func (e config) DatabaseUser() string {
	return getEnvironmentVariable("ROPC_DB_USER")
}

func (e config) DatabaseName() string {
	return getEnvironmentVariable("ROPC_DB_NAME")
}

func (e config) TokenSecret() string {
	return getEnvironmentVariable("ROPC_TOKEN_SECRET")
}

func (e config) DatabaseHost() string {
	return getEnvironmentVariable("ROPC_DB_HOST")
}

func (e config) DatabasePort() string {
	return getEnvironmentVariable("ROPC_DB_PORT")
}

func (e config) TokenExpiry() int {

	v, err := strconv.Atoi(getEnvironmentVariable("ROPC_TOKEN_EXPIRY"))
	if err != nil {
		NewLogger().Fatal("Error getting token expiry")
	}

	return v
}

func NewConfig() Config {
	return &config{}
}

func getEnvironmentVariable(env string) string {
	val, ok := os.LookupEnv(env)
	if !ok {
		log.Fatal("unable to load environment variable: " + env)
	}

	return val
}
