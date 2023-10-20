package services

import (
	"backend-server/mocks"
	"backend-server/model/entities"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CreateClientTest(t *testing.T) {
	t.Run("should panic if create client fails", func(t *testing.T) {
		clientRepository := new(mocks.ClientRepository)
		client := &entities.Application{ClientId: uuid.NewString()}
		clientRepository.On("Create", client).Return(errors.New("could not create client"))
		clientService := NewClientService(clientRepository)

		exec := func() {
			clientService.CreateClient(client)
		}

		assert.PanicsWithError(t, "could not create client", exec)
		clientRepository.AssertCalled(t, "Create", client)
	})

	t.Run("client id should not be empty", func(t *testing.T) {
		client := &entities.Application{ClientSecret: "secret"}
		clientRepository := new(mocks.ClientRepository)
		clientService := NewClientService(clientRepository)

		clientRepository.On("Create", client).Return(nil)

		exec := func() {
			clientService.CreateClient(client)
		}

		assert.PanicsWithError(t, "client id is required", exec)
		clientRepository.AssertNotCalled(t, "Create", client)
	})
}
