package conf

import (
	"backend-server/logger"
	"backend-server/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database[T any] interface {
	GetDatabaseConnection() *T
}

type database struct {
	conn *gorm.DB
}

func NewDataBase(config Config) Database[gorm.DB] {
	DbHost := config.DatabaseHost()
	DbUser := config.DatabaseUser()
	DbPassword := config.DatabasePassword()
	DbName := config.DatabaseName()
	DbPort := config.DatabasePort()

	DbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	db, err := gorm.Open(mysql.Open(DbUrl), &gorm.Config{})
	if err != nil {
		logger.Error(err.Error(), true)
	}

	err = db.AutoMigrate(&model.User{}, &model.Application{})
	if err != nil {
		logger.Error("Could not migrate:"+err.Error(), true)
	}

	logger.Info("Database connection established")

	return &database{db}
}

func (db *database) GetDatabaseConnection() *gorm.DB {
	return db.conn
}
