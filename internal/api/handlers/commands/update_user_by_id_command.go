package command

import (
	"context"

	cqrs "github.com/CallumLewisGH/Generic-Service-Base/internal/domain"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateUserByIdCommand(id uuid.UUID, user userModel.UpdateUserRequest) (*userModel.UserDTO, error) {
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
