package repos

import "gorm.io/gorm"

type BaseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) *BaseRepo {
	return &BaseRepo{db: db}
}

func (repo *BaseRepo) Model(model any) *BaseRepo {
	repo.db = repo.db.Model(model)
	return repo
}

func (repo *BaseRepo) Limit(limit int) *BaseRepo {
	repo.db = repo.db.Limit(limit)
	return repo
}

func (repo *BaseRepo) Offset(offset int) *BaseRepo {
	repo.db = repo.db.Offset(offset)
	return repo
}

func (repo *BaseRepo) Order(order string) *BaseRepo {
	repo.db = repo.db.Order(order)
	return repo
}

// Terminal Methods => Terminate the builder

// Terminal methods (execute the query)
func (r *BaseRepo) Find(dest any) error {
	return r.db.Find(dest).Error
}

func (r *BaseRepo) First(dest any) error {
	return r.db.First(dest).Error
}

func (repo *BaseRepo) Count() (int64, error) {
	var count int64
	err := repo.db.Count(&count).Error
	return count, err
}

func (repo *UserRepo) Delete() (int64, error) {
	result := repo.db.Delete(nil)
	return result.RowsAffected, result.Error
}
