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

func (repo *UserRepo) WithId(id uint) *UserRepo {
	repo.db = repo.db.Where("id = ?", id)
	return repo
}

func (repo *UserRepo) IsActive() *UserRepo {
	repo.db = repo.db.Where("is_active = ?", true)
	return repo
}

func (repo *UserRepo) IsInactive() *UserRepo {
	repo.db = repo.db.Where("is_active = ?", false)
	return repo
}
