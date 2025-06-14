package integration_tests

import (
	"testing"

	command "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/commands"
	query "github.com/CallumLewisGH/Generic-Service-Base/internal/api/handlers/queries"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestCreateUser_Success() {
	testUser := userModel.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		AuthId:   "test-auth-id",
	}

	createdUser, err := command.CreateUserCommand(testUser)
	s.NoError(err)
	s.NotNil(createdUser)
	s.Equal(testUser.Username, createdUser.Username)
	s.Equal(testUser.Email, createdUser.Email)
	s.Equal(testUser.AuthId, createdUser.AuthId)
}

func (s *TestSuite) TestGetAllUsers_Success() {
	// Create first test user
	testUser1 := userModel.CreateUserRequest{
		Username: "testuser1",
		Email:    "test1@example.com",
		AuthId:   "test-auth-id-1",
	}
	_, err := command.CreateUserCommand(testUser1)
	s.NoError(err)

	// Create second test user
	testUser2 := userModel.CreateUserRequest{
		Username: "testuser2",
		Email:    "test2@example.com",
		AuthId:   "test-auth-id-2",
	}
	_, err = command.CreateUserCommand(testUser2)
	s.NoError(err)

	// Get all users
	users, err := query.GetAllUsersQuery()
	s.NoError(err)
	s.GreaterOrEqual(len(users), 2)

	// Verify both test users are in the results
	var foundUser1, foundUser2 bool
	for _, user := range users {
		if user.Email == testUser1.Email {
			foundUser1 = true
		}
		if user.Email == testUser2.Email {
			foundUser2 = true
		}
	}
	s.True(foundUser1, "First test user should be in results")
	s.True(foundUser2, "Second test user should be in results")
}

func (s *TestSuite) TestGetUserById_Success() {
	// Create test user
	testUser := userModel.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		AuthId:   "test-auth-id",
	}
	createdUser, err := command.CreateUserCommand(testUser)
	s.NoError(err)

	// Get user by ID
	foundUser, err := query.GetUserByIdQuery(createdUser.ID)
	s.NoError(err)
	s.Equal(createdUser.ID, foundUser.ID)
	s.Equal(testUser.Username, foundUser.Username)
	s.Equal(testUser.Email, foundUser.Email)
}

func (s *TestSuite) TestGetUserById_NotFound() {
	nonExistentId := uuid.New()
	foundUser, err := query.GetUserByIdQuery(nonExistentId)
	s.Error(err)
	s.Contains(err.Error(), "record not found")
	s.Nil(foundUser)
}

func (s *TestSuite) TestUpdateUser_Success() {
	// Create test user
	testUser := userModel.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		AuthId:   "test-auth-id",
	}
	createdUser, err := command.CreateUserCommand(testUser)
	s.NoError(err)

	// Update user
	newUsername := "updated-username"
	newEmail := "updated@example.com"
	updateReq := userModel.UpdateUserRequest{
		Username: &newUsername,
		Email:    &newEmail,
	}
	updatedUser, err := command.UpdateUserByIdCommand(createdUser.ID, updateReq)
	s.NoError(err)
	s.Equal(newUsername, updatedUser.Username)
	s.Equal(newEmail, updatedUser.Email)
	s.Equal(createdUser.AuthId, updatedUser.AuthId) // AuthId should remain unchanged
}

func (s *TestSuite) TestUpdateUser_NotFound() {
	nonExistentId := uuid.New()
	newUsername := "updated-username"
	newEmail := "updated@example.com"
	updateReq := userModel.UpdateUserRequest{
		Username: &newUsername,
		Email:    &newEmail,
	}
	updatedUser, err := command.UpdateUserByIdCommand(nonExistentId, updateReq)
	s.Error(err)
	s.Contains(err.Error(), "record not found")
	s.Nil(updatedUser)
}

func (s *TestSuite) TestDeleteUser_Success() {
	// Create test user
	testUser := userModel.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		AuthId:   "test-auth-id",
	}
	createdUser, err := command.CreateUserCommand(testUser)
	s.NoError(err)

	// Delete user
	deletedUser, err := command.DeleteUserByIdCommand(createdUser.ID)
	s.NoError(err)
	s.Equal(createdUser.ID, deletedUser.ID)

	// Verify user is deleted
	_, err = query.GetUserByIdQuery(createdUser.ID)
	s.Error(err)
	s.Contains(err.Error(), "record not found")
}

func (s *TestSuite) TestDeleteUser_NotFound() {
	nonExistentId := uuid.New()
	deletedUser, err := command.DeleteUserByIdCommand(nonExistentId)
	s.Error(err)
	s.Contains(err.Error(), "record not found")
	s.Nil(deletedUser)
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
