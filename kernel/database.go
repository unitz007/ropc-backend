package kernel

import (
	"fmt"
	"ropc-backend/model"
	"ropc-backend/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	GetDatabaseConnection() any
}

type database struct {
	dbConn *gorm.DB
}

func (d database) GetDatabaseConnection() any {
	return d.dbConn
}

func NewDatabase(config utils.Config) (Database, error) {

	host := config.DatabaseHost()
	user := config.DatabaseUser()
	password := config.DatabasePassword()
	name := config.DatabaseName()

	DbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, name)

	db, err := gorm.Open(mysql.Open(DbUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&model.User{}, &model.Application{}, model.Test{})
	if err != nil {
		return nil, err
	}

	return database{
		dbConn: db,
	}, nil
}
