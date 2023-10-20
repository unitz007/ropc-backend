package conf

import (
	"backend-server/logger"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	configFileName = ".env"
	//configFileErrorMessage  = "Could not find config file: " + configFileName
	dbHost                  = "DB_HOST"
	dbPort                  = "DB_PORT"
	dbUser                  = "DB_USER"
	serverPort              = "SERVER_PORT"
	tokenSecret             = "TOKEN_SECRET"
	dbName                  = "DB_NAME"
	dbPassword              = "DB_PASSWORD"
	tokenExpiry             = "TOKEN_EXPIRY"
	tokenExpiryErrorMessage = "invalid token expiry"
)

//var (
//	config Config
//)

type Config interface {
	ServerPort() string
	DatabasePassword() string
	DatabaseUser() string
	DatabaseName() string
	TokenSecret() string
	DatabaseHost() string
	DatabasePort() string
	TokenExpiry() int
	Mux() string
}

type config struct{}

func (c *config) Mux() string {
	return getEnvironmentVariable("ROPC_MUX")
}

func NewConfig() Config {
	// load .env file
	err := godotenv.Load(configFileName)
	if err != nil {
		logger.Warn("Could not load .env file")
	}

	return &config{}
}

func (c *config) DatabaseName() string {
	return os.Getenv(dbName)
}

func (c *config) DatabasePort() string {
	return os.Getenv(dbPort)
}

func (c *config) ServerPort() string {
	return os.Getenv(serverPort)
}

func (c *config) DatabaseHost() string {
	return os.Getenv(dbHost)
}

func (c *config) DatabaseUser() string {
	return os.Getenv(dbUser)
}

func (c *config) DatabasePassword() string {
	return os.Getenv(dbPassword)
}

func (c *config) TokenSecret() string {
	return os.Getenv(tokenSecret)
}

func (c *config) TokenExpiry() int {
	return func() int {
		val, err := strconv.Atoi(os.Getenv(tokenExpiry))
		if err != nil {
			log.Fatal(tokenExpiryErrorMessage)
		}
		return val
	}()
}
