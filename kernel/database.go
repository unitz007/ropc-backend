package kernel

import (
	"fmt"
	"ropc-backend/model"
	"ropc-backend/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database[DatabaseConnectionReference any] interface {
	GetDatabaseConnection() *DatabaseConnectionReference
}

type database[DatabaseConnectionReference any] struct {
	dbConn *DatabaseConnectionReference
}

func (d database[DatabaseConnectionReference]) GetDatabaseConnection() *DatabaseConnectionReference {
	return d.dbConn
}

func NewDatabase(config utils.Config) (Database[gorm.DB], error) {

	host := config.DatabaseHost()
	user := config.DatabaseUser()
	password := config.DatabasePassword()
	name := config.DatabaseName()

	DbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, name)

	db, err := gorm.Open(mysql.Open(DbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//
	err = db.AutoMigrate(&model.User{}, &model.Application{})
	if err != nil {
		return nil, err
	}

	return database[gorm.DB]{
		dbConn: db,
	}, nil
}
