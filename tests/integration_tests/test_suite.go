package integration_tests

import (
	"github.com/CallumLewisGH/Generic-Service-Base/database"
	integration_test_config "github.com/CallumLewisGH/Generic-Service-Base/test_config"
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
	// Docker Container Setup
	suite.container = integration_test_config.StartTestDatabaseContainer(suite.T())

	// Database State Management
	database.SetTestMode(suite.container.DbConnStr)

	// Database Connection
	dbInstance := database.GetDatabase()
	suite.db = dbInstance.GetGormDatabase()

	// Migrations
	err := dbInstance.RunMigrations()
	suite.Require().NoError(err, "Failed to migrate database")
}

func (suite *TestSuite) TearDownSuite() {
	suite.T().Log("All tests completed - running cleanup")
	suite.container.Cleanup()
}

func (suite *TestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
}

func (suite *TestSuite) TearDownTest() {
	suite.tx.Rollback()
}
