package query

import (
	"context"

	cqrs "github.com/CallumLewisGH/Generic-Service-Base/internal/domain"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetUserByIdQuery(id uuid.UUID) (*userModel.UserDTO, error) {
	queryFunc := func(db *gorm.DB, ctx context.Context) (*userModel.UserDTO, error) {
		user := userModel.User{}
		user.ID = id
		if err := repos.NewUserRepo(db).First(&user); err != nil {
			return nil, err
		}

		return user.ToUserDTO(), nil
	}
	return cqrs.DbQuery(queryFunc)
}
