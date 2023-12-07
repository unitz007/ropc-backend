package kernel

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository[Model any] interface {
	Create(model Model) (*Model, error)
	GetById(id int) (*Model, error)
}

type repository[Model any] struct {
	model Model
	db    *gorm.DB
}

func (r repository[Model]) Create(model Model) (*Model, error) {
	fmt.Println("mode: ", r.model)
	fmt.Println("id:", r.db)
	err := r.db.Create(model).Error
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r repository[Model]) GetById(id int) (*Model, error) {
	//TODO implement me
	panic("implement me")
}

func NewRepository[Model any](model Model, db *gorm.DB) Repository[Model] {
	return repository[Model]{model: model, db: db}
}
