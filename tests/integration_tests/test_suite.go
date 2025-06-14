package integration_tests

import (
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
