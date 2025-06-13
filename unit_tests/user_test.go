package unitTests

import (
	"testing"

	"github.com/CallumLewisGH/Generic-Service-Base/database"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserRepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	tx   *gorm.DB
	repo *repos.UserRepo
}

func (s *UserRepoTestSuite) SetupSuite() {
	dbInstance := database.GetTestDatabase()
	s.db = dbInstance.GetGormDatabase()

	err := s.db.AutoMigrate(&userModel.User{})
	s.Require().NoError(err, "Failed to migrate database")
}

func (suite *UserRepoTestSuite) TearDownSuite() {
	suite.T().Log("All tests completed - running cleanup")
	err := suite.db.Migrator().DropTable(&userModel.User{})
	if err != nil {
		suite.T().Errorf("Failed to drop User table: %v", err)
		return
	}
	suite.T().Log("User table dropped successfully")

}

func (s *UserRepoTestSuite) SetupTest() {
	s.tx = s.db.Begin()

	adapter := database.NewGormAdapter(s.tx)
	s.repo = repos.NewUserRepo(adapter)
}

func (s *UserRepoTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *UserRepoTestSuite) TestGetUserByUsername_Success() {
	testUser := &userModel.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
	}
	err := s.tx.Create(testUser).Error
	s.Require().NoError(err)

	var foundUser userModel.User
	err = s.repo.WithName("testuser").First(&foundUser)

	s.NoError(err)
	s.Equal(testUser.ID, foundUser.ID)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
