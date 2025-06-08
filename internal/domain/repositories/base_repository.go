package repos

import "gorm.io/gorm"

type BaseRepo[T any] struct {
	db *gorm.DB
}

func NewBaseRepo[T any](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{db: db}
}

func (repo *BaseRepo[T]) Limit(limit int) *BaseRepo[T] {
	repo.db = repo.db.Limit(limit)
	return repo
}

func (repo *BaseRepo[T]) Offset(offset int) *BaseRepo[T] {
	repo.db = repo.db.Offset(offset)
	return repo
}

func (repo *BaseRepo[T]) Order(order string) *BaseRepo[T] {
	repo.db = repo.db.Order(order)
	return repo
}

// Terminal Methods => Terminate the builder
func (repo *BaseRepo[T]) Find(dest *[]T) error {
	return repo.db.Find(dest).Error
}

func (repo *BaseRepo[T]) First(dest *T) error {
	return repo.db.First(dest).Error
}

func (repo *BaseRepo[T]) Count() (int64, error) {
	var count int64
	err := repo.db.Count(&count).Error
	return count, err
}

func (repo *BaseRepo[T]) Create(dest *T) error {
	return repo.db.Create(dest).Error
}

func (repo *BaseRepo[T]) Delete(dest *T) (int64, error) {
	result := repo.db.Delete(dest)
	return result.RowsAffected, result.Error
}
