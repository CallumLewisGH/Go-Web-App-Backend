package query

import (
	"context"

	cqrs "github.com/CallumLewisGH/Generic-Service-Base/internal/domain"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/domain/models"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	"gorm.io/gorm"
)

func GetAllUsersQuery() ([]models.User, error) {
	queryFunc := func(db *gorm.DB, ctx context.Context) ([]models.User, error) {
		users := []models.User{}
		usersRepo := repos.NewUserRepo(db)
		err := usersRepo.Find(&users)

		if err != nil {
			return nil, err
		}
		return users, nil
	}
	return cqrs.DbQuery(queryFunc)
}
