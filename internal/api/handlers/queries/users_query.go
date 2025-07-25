package query

import (
	"context"

	cqrs "github.com/CallumLewisGH/Generic-Service-Base/internal/domain"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"gorm.io/gorm"
)

func GetAllUsersQuery() ([]userModel.UserDTO, error) {
	queryFunc := func(db *gorm.DB, ctx context.Context) ([]userModel.UserDTO, error) {
		users := []userModel.User{}
		if err := repos.NewUserRepo(db).Find(&users); err != nil {
			return nil, err
		}

		return userModel.ToUserDTOs(users), nil
	}
	return cqrs.DbQuery(queryFunc)
}
