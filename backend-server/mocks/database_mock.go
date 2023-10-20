package mocks

import (
	"backend-server/model"
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	mysqltestcontainer "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseMock struct {
	db *gorm.DB
}

type Container struct {
	Container *mysqltestcontainer.MySQLContainer
	URI       string
}

type DBMock struct {
	Mock sqlmock.Sqlmock
	db   *gorm.DB
}

func NewDBMock(t testing.TB) *DBMock {
	t.Helper()
	mockDb, m, err := sqlmock.New()
	require.NoError(t, err)

	dialect := mysql.New(mysql.Config{
		Conn:       mockDb,
		DriverName: "mysql",
	})

	db, _ := gorm.Open(dialect, &gorm.Config{})
	//require.NoError(t, err)

	return &DBMock{
		Mock: m,
		db:   db,
	}
}

func (m *DBMock) GetDatabaseConnection() *gorm.DB {
	return m.db
}

func prepareContainer(ctx context.Context) *Container {
	dbContainer, err := mysqltestcontainer.RunContainer(ctx,
		mysqltestcontainer.WithDatabase("ropc_db"),
		mysqltestcontainer.WithPassword("testpassword"),
		mysqltestcontainer.WithUsername("testusername"),
	)

	if err != nil {
		panic(err)
	}

	hostIP, err := dbContainer.Host(ctx)
	if err != nil {
		panic(err)
	}

	mappedPort, err := dbContainer.MappedPort(ctx, "3306")
	if err != nil {
		panic(err)
	}

	return &Container{
		Container: dbContainer,
		URI:       fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "testusername", "testpassword", hostIP, mappedPort.Port(), "ropc_db"),
	}
}

func NewDatabaseMock(t testing.TB) *DatabaseMock {
	t.Helper()
	ctx := context.Background()

	c := prepareContainer(ctx)

	db, err := gorm.Open(mysql.Open(c.URI), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.AutoMigrate(&model.Application{}, &model.User{})
	if err != nil {
		log.Fatal("Could not migrate:", err.Error())
	}

	defer t.Cleanup(func() {
		err := c.Container.Terminate(ctx)
		if err != nil {
			log.Fatal(err)
		}
	})

	return &DatabaseMock{
		db: db,
	}
}

func (m *DatabaseMock) GetDatabaseConnection() *gorm.DB {
	return m.db
}
