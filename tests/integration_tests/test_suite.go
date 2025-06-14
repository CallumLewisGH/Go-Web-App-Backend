package integration_tests

import (
	"github.com/CallumLewisGH/Generic-Service-Base/database"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	database  *database.Database
	container *DockerContainer
}

func (suite *TestSuite) SetupSuite() {
	// Docker Container Setup
	suite.container = StartTestDatabaseContainer(suite.T())

	// Database State Management
	database.SetTestMode(suite.container.DbConnStr)

	// Database Connection
	suite.database = database.GetDatabase()

	// Migrations
	err := suite.database.RunMigrations()
	suite.Require().NoError(err, "Failed to migrate database")
}

func (suite *TestSuite) TearDownSuite() {
	suite.T().Log("All tests completed - running cleanup")
	suite.container.Cleanup()
}

func (suite *TestSuite) TearDownTest() {
	suite.database.ClearAllTables()
}
