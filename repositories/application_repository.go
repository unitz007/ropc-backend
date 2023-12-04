package repositories

import (
	"errors"
	"ropc-backend/kernel"
	"ropc-backend/model"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	GetByClientId(clientId string) (*model.Application, error)
	GetByClientIdAndUserId(clientId string, userId uint) (*model.Application, error)
	GetAll(userId uint) []model.Application
	Create(client *model.Application) error
	Update(app *model.Application) (*model.Application, error)
	GetByName(name string) (*model.Application, error)
	Delete(id uint) error
	GetByNameAndUserId(name string, userId uint) (*model.Application, error)
}

type applicationRepository struct {
	db kernel.Database
}

func (a applicationRepository) GetByNameAndUserId(name string, userId uint) (*model.Application, error) {

	client, err := a.db.DatabaseConnection().CustomSingleQuery("name = ? and user_id = ?", name, userId)

	if err != nil {
		return nil, errors.New("application not found")
	}

	return client.(*model.Application), nil
}

func (a applicationRepository) Delete(id uint) error {
	return a.db.GetDatabaseConnection().(*gorm.DB).
		Unscoped().
		Model(&model.Application{}).
		Delete("id = ?", id).
		Error
}

func (a applicationRepository) Update(app *model.Application) (*model.Application, error) {
	err := a.db.GetDatabaseConnection().(*gorm.DB).
		Model(app).
		Where("client_id = ?", app.ClientId).
		Update("client_secret", app.ClientSecret).
		Error

	if err != nil {
		return nil, err
	}

	return app, nil
}

func NewApplicationRepository(db kernel.Database) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (a applicationRepository) GetByClientIdAndUserId(clientId string, userId uint) (*model.Application, error) {

	var client model.Application

	err := a.db.GetDatabaseConnection().(*gorm.DB).
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

	err := a.db.DatabaseConnection().(*gorm.DB).
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

	err := a.db.GetDatabaseConnection().(*gorm.DB).
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

	a.db.GetDatabaseConnection().(*gorm.DB).
		Model(model.Application{}).
		Where("user_id = ?", userId).
		Scan(&clients)

	return clients
}

func (a applicationRepository) Create(client *model.Application) error {
	client.ClientId = uuid.NewString()

	err := a.db.GetDatabaseConnection().(*gorm.DB).Create(client).Error

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
