package integration_tests

import (
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

// Add any unit test specific setup logic here
// I don't think I'll need it but it's nice to have
func (suite *TestSuite) SetupSuite() {
}

func (suite *TestSuite) TearDownSuite() {
}

func (suite *TestSuite) SetupTest() {
}

func (suite *TestSuite) TearDownTest() {
}
