package integration_tests

import (
	"testing"

	command "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/commands"
	query "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/queries"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"github.com/stretchr/testify/suite"
)

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestGetUserByAuthId_Success() {
	testUser := userModel.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		AuthId:   "test-auth-id",
	}

	_, err := command.CreateUserCommand(testUser)
	s.NoError(err)

	outTestUser, err := query.GetUserByAuthIdQuery(testUser.AuthId)

	s.NoError(err)
	s.Equal(outTestUser.Email, testUser.Email)
}

func (s *TestSuite) TestGetUserByAuthId_Fail() {
	nonExistentAuthId := "non-existent-auth-id"
	outTestUser, err := query.GetUserByAuthIdQuery(nonExistentAuthId)

	s.Contains(err.Error(), "record not found")
	s.Nil(outTestUser)
}
