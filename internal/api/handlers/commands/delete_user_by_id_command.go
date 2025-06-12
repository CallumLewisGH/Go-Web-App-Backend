package command

import (
	"context"

	"github.com/CallumLewisGH/Generic-Service-Base/database"
	cqrs "github.com/CallumLewisGH/Generic-Service-Base/internal/domain"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"github.com/google/uuid"
)

func DeleteUserByIdCommand(id uuid.UUID) (*userModel.UserDTO, error) {
	commandFunc := func(db database.IDatabase, ctx context.Context) (*userModel.UserDTO, error) {
		user := userModel.User{}
		user.ID = id
		delted_user, err := repos.NewUserRepo(db).DeleteOne(user)

		if err != nil {
			return nil, err
		}

		return delted_user.ToUserDTO(), nil
	}
	return cqrs.DbExecute(commandFunc)
}
