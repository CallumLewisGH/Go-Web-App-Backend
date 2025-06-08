package repos

import (
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"gorm.io/gorm"
)

type UserRepo struct {
	*BaseRepo[userModel.User]
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{BaseRepo: NewBaseRepo[userModel.User](db)}
}

func (repo *UserRepo) WithName(name string) *UserRepo {
	repo.db = repo.db.Where("name = ?", name)
	return repo
}

func (repo *UserRepo) WithId(ID int) *UserRepo {
	repo.db = repo.db.Where("ID = ?", ID)
	return repo
}
