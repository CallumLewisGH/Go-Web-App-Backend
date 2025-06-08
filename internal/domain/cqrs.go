package cqrs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CallumLewisGH/Generic-Service-Base/database"
	"gorm.io/gorm"
)

func DbQuery[T any](queryFunction func(*gorm.DB, context.Context) ([]T, error)) ([]T, error) {
	db := database.GetDatabase().GetGormDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx := db.WithContext(ctx)

	return queryFunction(tx, ctx)
}

func DbExecute[T any](commandFunction func(*gorm.DB, context.Context) ([]T, error)) ([]T, error) {
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

	result, err := commandFunction(tx, ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("operation failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return result, nil
}
