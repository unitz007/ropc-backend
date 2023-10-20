package services

import (
	"backend-server/model"
	"backend-server/repositories"
	"errors"
)

type ApplicationService interface {
	Get() []model.Application
	CreateClient(client *model.Application)
}

type clientService struct {
	clientRepository repositories.ApplicationRepository
}

func (s *clientService) CreateClient(client *model.Application) {

	if client.ClientId == "" {
		panic(errors.New("client id is required"))
	}

	err := s.clientRepository.Create(client)

	if err != nil {
		panic(err)
	}
}

func NewClientService(repository repositories.ApplicationRepository) ApplicationService {
	return &clientService{clientRepository: repository}
}

func (s *clientService) Get() []model.Application {
	return s.clientRepository.GetAll()
}
