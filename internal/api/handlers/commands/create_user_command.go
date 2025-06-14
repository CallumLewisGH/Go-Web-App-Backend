package command

import (
	"context"

	cqrs "github.com/CallumLewisGH/Generic-Service-Base/internal/domain"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"gorm.io/gorm"
)

func CreateUserCommand(userReq userModel.CreateUserRequest) (*userModel.UserDTO, error) {
	commandFunc := func(db *gorm.DB, ctx context.Context) (*userModel.UserDTO, error) {
		user := userModel.User{
			Username: userReq.Username,
			Email:    userReq.Email,
			AuthId:   userReq.AuthId,
		}

		created_user, err := repos.NewUserRepo(db).CreateOne(user)

		if err != nil {
			return nil, err
		}

		return created_user.ToUserDTO(), nil
	}

	return cqrs.DbExecute(commandFunc)

}
