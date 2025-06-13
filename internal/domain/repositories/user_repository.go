package repos

import (
	"github.com/CallumLewisGH/Generic-Service-Base/database"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"github.com/google/uuid"
)

type UserRepo struct {
	*BaseRepo[userModel.User]
}

func NewUserRepo(db database.IDatabase) *UserRepo {
	return &UserRepo{BaseRepo: NewBaseRepo[userModel.User](db)}
}

func (repo *UserRepo) WithName(name string) *UserRepo {
	repo.db = repo.db.Where("username = ?", name)
	return repo
}

func (repo *UserRepo) WithId(id uuid.UUID) *UserRepo {
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
