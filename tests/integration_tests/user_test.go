package integration_tests

import (
	"context"
	"testing"

	cqrs "github.com/CallumLewisGH/Generic-Service-Base/internal/domain"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	integration_test_config "github.com/CallumLewisGH/Generic-Service-Base/test_config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	db        *gorm.DB
	tx        *gorm.DB
	container *integration_test_config.DockerContainer
}

func (suite *TestSuite) SetupSuite() {
	suite.T().Log("Spinning up test db container")

	// Docker Container Setup

	suite.container = integration_test_config.StartTestDatabaseContainer(suite.T())

	// Database Connection
	dbInstance := integration_test_config.GetTestDatabase(suite.container.DbConnStr)

	suite.db = dbInstance.GetTestGormDatabase()

	// Migrations
	err := dbInstance.RunTestMigrations()
	suite.Require().NoError(err, "Failed to migrate database")
}

func (s *TestSuite) TearDownSuite() {
	s.T().Log("All tests completed - running cleanup")
	s.container.Cleanup()
}

func (s *TestSuite) SetupTest() {
	s.tx = s.db.Begin()
}

func (s *TestSuite) TearDownTest() {
	s.tx.Rollback()
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestGetUserByUsername_Success(t *testing.T) {
	testUser := userModel.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
	}
	var foundUser userModel.User

	handlerFunc := func(db *gorm.DB, ctx context.Context) (*userModel.UserDTO, error) {
		testUser, err := repos.NewUserRepo(s.tx).CreateOne(testUser)
		s.Require().NoError(err)

		return testUser.ToUserDTO(), nil
	}

	outTestUser, err := cqrs.TestDbExecute(handlerFunc, s.container.DbConnStr)
	s.NoError(err)

	handlerFunc = func(db *gorm.DB, ctx context.Context) (*userModel.UserDTO, error) {
		err := repos.NewUserRepo(s.tx).WithName("testuser").First(&foundUser)
		s.Require().NoError(err)

		return foundUser.ToUserDTO(), nil
	}

	outFoundUser, err := cqrs.TestDbQuery(handlerFunc, s.container.DbConnStr)

	s.NoError(err)
	s.Equal(outTestUser.ID, outFoundUser.ID)
}
