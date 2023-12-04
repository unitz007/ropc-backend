package kernel

import (
	"fmt"
	"ropc-backend/model"
	"ropc-backend/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database interface {
	DatabaseConnection() DataBaseConnection[any]
}

type DataBaseConnection[Model any] interface {
	Create(entity Model) (*Model, error)
	GetById(id int) (*Model, error)
	GetList() ([]Model, error)
	CustomSingleQuery(query string, vals ...any) (*Model, error)
	CustomListQuery(query string, vals ...any) ([]Model, error)
	Delete(id int) error
}

type database[Model any] struct {
	dbConn *gorm.DB
}

func (d database[Model]) Delete(id int) error {

	var m Model

	return d.dbConn.Unscoped().
		Model(m).
		Delete("id = ?", id).
		Error
}

func (d database[Model]) CustomSingleQuery(query string, cond ...any) (*Model, error) {

	var m Model

	err := d.dbConn.
		Model(m).
		Where(query, cond).
		First(&m).
		Error

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (d database[Model]) CustomListQuery(query string, cond ...any) ([]Model, error) {
	//TODO implement me
	panic("implement me")
}

func (d database[Model]) Create(m Model) (*Model, error) {
	if err := d.dbConn.Create(m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

func (d database[Model]) GetById(id int) (*Model, error) {

	var res Model
	err := d.dbConn.Where("id = ?", id).First(res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (d database[Model]) GetList() ([]Model, error) {
	//TODO implement me
	panic("implement me")
}

func (d database[Model]) DatabaseConnection() *gorm.DB {
	return d.dbConn
}

func NewDatabase[Model any](config utils.Config) (DataBaseConnection[Model], error) {

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

	return database[Model]{
		dbConn: db,
	}, nil
}
