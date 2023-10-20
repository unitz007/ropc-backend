package repositories

import (
	"backend-server/conf"
	"backend-server/model"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	GetByClientId(clientId string) (*model.Application, error)
	GetAll() []model.Application
	Create(client *model.Application) error
	Update(app *model.Application) (*model.Application, error)
	GetByName(name string) (*model.Application, error)
}

type applicationRepository struct {
	db conf.Database[gorm.DB]
}

func (a applicationRepository) Update(app *model.Application) (*model.Application, error) {
	err := a.db.GetDatabaseConnection().
		Model(app).
		Where("client_id = ?", app.ClientId).
		Update("client_secret", app.ClientSecret).
		Error

	if err != nil {
		return nil, err
	}

	return app, nil
}

func NewApplicationRepository(db conf.Database[gorm.DB]) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (a applicationRepository) GetByClientId(clientId string) (*model.Application, error) {

	var client model.Application

	err := a.db.GetDatabaseConnection().
		Model(&model.Application{}).
		Where("client_id = ?", clientId).
		First(&client).
		Error

	if err != nil {
		return nil, errors.New("application not found")
	}

	return &client, nil
}

func (a applicationRepository) GetByName(name string) (*model.Application, error) {

	var client model.Application

	err := a.db.GetDatabaseConnection().
		Model(&model.Application{}).
		Where("name = ?", name).
		First(&client).
		Error

	if err != nil {
		return nil, errors.New("application not found")
	}

	return &client, nil
}

func (a applicationRepository) GetAll() []model.Application {

	var clients []model.Application

	a.db.GetDatabaseConnection().Find(&clients)

	return clients
}

func (a applicationRepository) Create(client *model.Application) error {

	err := a.db.GetDatabaseConnection().Create(client).Error

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.New("application already exists")
		} else {
			return errors.New("could not create application. Contact administrator")
		}
	}

	return err
}
