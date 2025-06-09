package command

import (
	"context"

	cqrs "github.com/CallumLewisGH/Generic-Service-Base/internal/domain"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"gorm.io/gorm"
)

func UpdateUserByIdCommand(id uint, user userModel.UserDTO) (*userModel.UserDTO, error) {
	commandFunc := func(db *gorm.DB, ctx context.Context) (*userModel.UserDTO, error) {
		dest := userModel.User{}
		dest.ID = id
		updated_user, err := repos.NewUserRepo(db).UpdateOne(dest, user)

		if err != nil {
			return nil, err
		}

		return updated_user.ToUserDTO(), nil
	}
	return cqrs.DbExecute(commandFunc)
}
