package cqrs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CallumLewisGH/Generic-Service-Base/database"
	"gorm.io/gorm"
)

// DbQuery runs a read-only query and automatically infers whether the result is []T or *T.
func DbQuery[R any](queryFunc func(*gorm.DB, context.Context) (R, error)) (R, error) {
	db := database.GetDatabase().GetGormDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx := db.WithContext(ctx)
	return queryFunc(tx, ctx)
}

// DbExecute runs a write operation (insert/update/delete) and infers []T or *T.
func DbExecute[R any](commandFunc func(*gorm.DB, context.Context) (R, error)) (R, error) {
	db := database.GetDatabase().GetGormDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx := db.Begin().WithContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic in Execute: %v", r)
		}
	}()

	result, err := commandFunc(tx, ctx)
	if err != nil {
		tx.Rollback()
		return result, fmt.Errorf("operation failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return result, fmt.Errorf("commit failed: %w", err)
	}

	return result, nil
}
