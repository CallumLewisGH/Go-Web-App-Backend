package repos

import (
	"gorm.io/gorm"
)

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

// Find all items that match the dest
func (repo *BaseRepo[T]) Find(dest *[]T) error {
	return repo.db.Find(dest).Error
}

// Gets the first item that match the dest
func (repo *BaseRepo[T]) First(dest *T) error {
	return repo.db.First(dest).Error
}

// Count counts all items that match the dest conditions
func (repo *BaseRepo[T]) Count(dest T) (int64, error) {
	var count int64
	err := repo.db.Model(dest).Where(dest).Count(&count).Error
	return count, err
}

// CreateOne accepts and returns values
func (repo *BaseRepo[T]) CreateOne(dest T) (T, error) {
	if err := repo.db.Create(&dest).Error; err != nil {
		return *new(T), err
	}
	return dest, nil
}

// CreateMany accepts and returns slices
func (repo *BaseRepo[T]) CreateMany(dest []T) ([]T, error) {
	if err := repo.db.Create(&dest).Error; err != nil {
		return nil, err
	}
	return dest, nil
}

// DeleteOne returns the deleted record
func (repo *BaseRepo[T]) DeleteOne(dest T) (T, error) {
	var model T
	if err := repo.db.Where(dest).First(&model).Error; err != nil {
		return *new(T), err
	}
	if err := repo.db.Where(dest).Delete(&dest).Error; err != nil {
		return *new(T), err
	}
	return model, nil
}

// DeleteMany returns deleted records
func (repo *BaseRepo[T]) DeleteMany(dest []T) ([]T, error) {
	var toDelete []T
	if err := repo.db.Where(dest).Find(&toDelete).Error; err != nil {
		return nil, err
	}
	if err := repo.db.Where(dest).Delete(&dest).Error; err != nil {
		return nil, err
	}
	return toDelete, nil
}

// UpdateOne returns updated record
func (repo *BaseRepo[T]) UpdateOne(dest T, updates any) (T, error) {
	if err := repo.db.Model(&dest).Updates(updates).Error; err != nil {
		return *new(T), err
	}
	var updated T
	if err := repo.db.Where(&dest).First(&updated).Error; err != nil {
		return *new(T), err
	}
	return updated, nil
}

// UpdateMany returns updated records
func (repo *BaseRepo[T]) UpdateMany(conditions T, updates any) ([]T, error) {
	var models []T
	if err := repo.db.Where(conditions).Find(&models).Error; err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return models, nil
	}
	if err := repo.db.Model(new(T)).Where(conditions).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := repo.db.Where(conditions).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}
