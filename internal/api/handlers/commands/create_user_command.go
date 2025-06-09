package command

import (
	"context"

	"github.com/CallumLewisGH/Generic-Service-Base/internal/api/authentication"
	cqrs "github.com/CallumLewisGH/Generic-Service-Base/internal/domain"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"gorm.io/gorm"
)

func CreateUserCommand(userReq *userModel.UserRequest) (*userModel.UserDTO, error) {
	commandFunc := func(db *gorm.DB, ctx context.Context) (*userModel.UserDTO, error) {
		generatedHash, _ := authentication.GenerateHash(userReq.Password)

		user := userModel.User{
			Username:     userReq.Username,
			Email:        userReq.Email,
			PasswordHash: generatedHash,
		}

		created_user, err := repos.NewUserRepo(db).CreateOne(user)

		if err != nil {
			return nil, err
		}

		return created_user.ToUserDTO(), nil
	}

	return cqrs.DbExecute(commandFunc)

}
