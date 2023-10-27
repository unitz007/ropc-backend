package repositories

import (
	"backend-server/model"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	GetByClientId(clientId string) (*model.Application, error)
	GetByClientIdAndUserId(clientId string, userId uint) (*model.Application, error)
	GetAll(userId uint) []model.Application
	Create(client *model.Application) error
	Update(app *model.Application) (*model.Application, error)
	GetByName(name string) (*model.Application, error)
}

type applicationRepository struct {
	db Database[gorm.DB]
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

func NewApplicationRepository(db Database[gorm.DB]) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (a applicationRepository) GetByClientIdAndUserId(clientId string, userId uint) (*model.Application, error) {

	var client model.Application

	err := a.db.GetDatabaseConnection().
		Model(&model.Application{}).
		Where("client_id = ? AND user_id = ?", clientId, userId).
		First(&client).
		Error

	if err != nil {
		return nil, errors.New("application not found")
	}

	return &client, nil
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

func (a applicationRepository) GetAll(userId uint) []model.Application {

	var clients []model.Application

	a.db.GetDatabaseConnection().
		Raw("select * from applications where user_id = ?", userId).
		Scan(&clients)

	return clients
}

func (a applicationRepository) Create(client *model.Application) error {

	err := a.db.GetDatabaseConnection().Create(client).Error

	if err != nil {
		if strings.Contains(err.Error(), "name") {
			return errors.New("application name already exists")
		} else if strings.Contains(err.Error(), "client_id") {
			return errors.New("client_id already exists")
		} else {
			return errors.New("could not create application. Contact administrator")
		}
	}

	return err
}
