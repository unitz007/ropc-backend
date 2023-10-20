package repositories

import (
	"backend-server/conf"
	"backend-server/mocks"
	"backend-server/model/entities"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateClient(t *testing.T) {
	t.Skip()
	var db conf.Database[gorm.DB]

	t.Run("should return error if client_id already exists", func(t *testing.T) {
		db = mocks.NewDatabaseMock(t)
		repo := NewApplicationRepository(db)

		clientId := "client_id"
		client := &entities.Application{
			ClientId:     clientId,
			ClientSecret: mock.Anything,
		}

		err := repo.Create(client)
		require.NoError(t, err, "should not return error")

		err = repo.Create(client)
		require.Errorf(t, err, "should return error")

		expectedErrorMsg := "client Id already exists"

		if err.Error() != expectedErrorMsg {
			t.Errorf("expected %s, got %s", expectedErrorMsg, err.Error())
		}
	})

	t.Run("should fail if client secret or client id is empty", func(t *testing.T) {
		db = mocks.NewDatabaseMock(t)
		repository := NewApplicationRepository(db)
		err := repository.Create(&entities.Application{ClientId: "clientId"})
		require.Errorf(t, err, "should return error cause client secret is empty")
	})

}
