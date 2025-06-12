package cqrs

import (
	"context"
	"testing"
	"time"

	"github.com/CallumLewisGH/Generic-Service-Base/database"
)

// MockQuery runs a read-only query against a mock database
func MockQuery[R any](
	t *testing.T,
	mockDB database.IDatabase,
	queryFunc func(database.IDatabase, context.Context) (R, error),
) (R, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx := mockDB.WithContext(ctx)
	return queryFunc(tx, ctx)
}

// MockExecute runs a write operation against a mock database
func MockExecute[R any](
	t *testing.T,
	mockDB database.IDatabase,
	commandFunc func(database.IDatabase, context.Context) (R, error),
) (R, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx := mockDB.Begin()
	ctxTx := tx.WithContext(ctx)

	result, err := commandFunc(ctxTx, ctx)
	if err != nil {
		tx.Rollback()
		t.Errorf("Operation failed: %v", err)
		return result, err
	}

	if err := tx.Commit().Error(); err != nil {
		t.Errorf("Commit failed: %v", err)
		return result, err
	}

	return result, nil
}
