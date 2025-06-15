package integration_tests

import (
	"github.com/CallumLewisGH/Generic-Service-Base/database"
	"github.com/CallumLewisGH/Generic-Service-Base/tests/test_config"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	gormDatabase      *database.Database
	redisDatabase     *database.Redis
	databaseContainer *test_config.DockerContainer
	redisContainer    *test_config.DockerContainer
}

func (suite *TestSuite) SetupSuite() {
	// Database Docker Container Setup
	suite.databaseContainer = test_config.StartTestDatabaseContainer(suite.T())

	// Database State Management
	database.SetTestMode(suite.databaseContainer.DbConnStr)

	// Database Connection
	suite.gormDatabase = database.GetDatabase()

	// Migrations
	err := suite.gormDatabase.RunMigrations()
	suite.Require().NoError(err, "Failed to migrate database")

	// Redis Docker Container Setup
	suite.redisContainer = test_config.StartTestRedisContainer(suite.T())

	// Redis Connection
	suite.redisDatabase = database.GetRedis()
}

func (suite *TestSuite) TearDownSuite() {
	suite.T().Log("All tests completed - running cleanup")
	suite.databaseContainer.Cleanup()
	suite.redisContainer.Cleanup()
}

func (suite *TestSuite) TearDownTest() {
	suite.gormDatabase.ClearAllTables()
	suite.redisDatabase.ClearAll()
}
