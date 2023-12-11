package kernel

import (
	"gorm.io/gorm"
)

type Repository[Model any] interface {
	Create(model Model) error
	Get(conditions string) (*Model, error)
	GetAll(conditions string) []Model
	Delete(conditions string) error
	Update(id uint, fields map[string]any) error
}

type repository[Model any] struct {
	model Model
	db    *gorm.DB
}

func (r repository[Model]) GetAll(conditions string) []Model {
	var models []Model
	r.db.
		Model(r.model).
		Where(conditions).
		Scan(&models)

	return models
}

func (r repository[Model]) Get(conditions string) (*Model, error) {
	var m Model
	err := r.db.Unscoped().
		Where(conditions).
		First(&m).
		Error

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r repository[Model]) Create(model Model) error {
	err := r.db.Create(&model).Error
	if err != nil {
		return err
	}

	return nil
}

func (r repository[Model]) Delete(condition string) error {
	return r.db.
		Unscoped().
		Delete(&r.model, condition).
		Error
}

func (r repository[Model]) Update(id uint, fields map[string]any) error {
	err := r.db.
		Model(&r.model).
		Where("id = ?", id).
		UpdateColumns(fields).
		Error

	if err != nil {
		return err
	}

	return nil
}

func NewRepository[Model any](model Model, db *gorm.DB) Repository[Model] {
	return repository[Model]{model: model, db: db}
}
