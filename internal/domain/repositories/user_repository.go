package repos

import (
	"github.com/CallumLewisGH/Generic-Service-Base/internal/domain/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	*BaseRepo
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{BaseRepo: NewBaseRepo(db.Model(&models.User{}))}
}

func (repo *UserRepo) WithName(name string) *UserRepo {
	repo.db = repo.db.Where("name = ?", name)
	return repo
}

func (repo *UserRepo) WithId(ID int) *UserRepo {
	repo.db = repo.db.Where("ID = ?", ID)
	return repo
}
