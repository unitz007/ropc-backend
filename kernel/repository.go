package kernel

import (
	"gorm.io/gorm"
)

type Repository[Model any] interface {
	Create(model *Model) (*Model, error)
	GetById(id int) (*Model, error)
	QuerySingleResult(model Model) (*Model, error)
	Delete(id uint) error
}

type repository[Model any] struct {
	model Model
	db    *gorm.DB
}

func (r repository[Model]) QuerySingleResult(qModel Model) (*Model, error) {
	//var model Model
	err := r.db.
		////Model(r.model).
		//Where(query, condition).
		First(&qModel).
		Error

	if err != nil {
		return nil, err
	}

	return &qModel, nil
}

func (r repository[Model]) Create(model *Model) (*Model, error) {
	err := r.db.Create(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r repository[Model]) GetById(id int) (*Model, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository[Model]) Delete(id uint) error {
	return r.db.
		Unscoped().
		Model(r.model).
		Delete("where id = ?", id).
		Error
}

func NewRepository[Model any](model Model, db *gorm.DB) Repository[Model] {
	return repository[Model]{model: model, db: db}
}
