package cqrs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CallumLewisGH/Generic-Service-Base/database"
)

// DbQuery runs a read-only query using IDatabase
func DbQuery[R any](queryFunc func(database.IDatabase, context.Context) (R, error)) (R, error) {
	db := database.GetDatabase().GetGormDatabase()
	adapter := database.NewGormAdapter(db)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx := adapter.WithContext(ctx)
	return queryFunc(tx, ctx)
}

// DbExecute runs a write operation using IDatabase
func DbExecute[R any](commandFunc func(database.IDatabase, context.Context) (R, error)) (R, error) {
	db := database.GetDatabase().GetGormDatabase()
	adapter := database.NewGormAdapter(db)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx := adapter.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic in Execute: %v", r)
		}
	}()

	result, err := commandFunc(tx.WithContext(ctx), ctx)
	if err != nil {
		tx.Rollback()
		return result, fmt.Errorf("operation failed: %w", err)
	}

	if err := tx.Commit().Error(); err != nil {
		return result, fmt.Errorf("commit failed: %w", err)
	}

	return result, nil
}
