package conf

import (
	"backend-server/logger"
	"os"
	"strconv"
)

type envConfig struct{}

func (e envConfig) Mux() string {
	return getEnvironmentVariable("ROPC_MUX")
}

func (e envConfig) ServerPort() string {
	return getEnvironmentVariable("ROPC_SERVER_PORT")
}

func (e envConfig) DatabasePassword() string {
	return getEnvironmentVariable("ROPC_DATABASE_PASSWORD")
}

func (e envConfig) DatabaseUser() string {
	return getEnvironmentVariable("ROPC_DB_USER")
}

func (e envConfig) DatabaseName() string {
	return getEnvironmentVariable("ROPC_DB_NAME")
}

func (e envConfig) TokenSecret() string {
	return getEnvironmentVariable("ROPC_TOKEN_SECRET")
}

func (e envConfig) DatabaseHost() string {
	return getEnvironmentVariable("ROPC_DB_HOST")
}

func (e envConfig) DatabasePort() string {
	return getEnvironmentVariable("ROPC_DB_PORT")
}

func (e envConfig) TokenExpiry() int {

	v, err := strconv.Atoi(getEnvironmentVariable("ROPC_TOKEN_EXPIRY"))
	if err != nil {
		logger.Error("Error getting token expiry", true)
	}

	return v
}

func NewEnvConfig() Config {
	return &envConfig{}
}

func getEnvironmentVariable(env string) string {
	val, ok := os.LookupEnv(env)
	if !ok {
		logger.Error("unable to load environment variable: "+env, true)
	}

	return val
}
