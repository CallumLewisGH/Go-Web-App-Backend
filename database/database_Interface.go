package database

import (
	"context"
	"database/sql"
)

type IDatabase interface {
	Where(query any, args ...any) IDatabase
	First(dest any, conds ...any) IDatabase
	Create(value any) IDatabase
	Delete(value any, conds ...any) IDatabase
	Model(value any) IDatabase
	Updates(values any) IDatabase
	Limit(limit int) IDatabase
	Offset(offset int) IDatabase
	Order(value any) IDatabase
	Count(count *int64) IDatabase
	Find(dest any, conds ...any) IDatabase
	Select(query any, args ...any) IDatabase

	WithContext(ctx context.Context) IDatabase
	Begin(opts ...*sql.TxOptions) IDatabase
	Commit() IDatabase
	Rollback() IDatabase

	DB() (*sql.DB, error)
	AutoMigrate(dst ...any) error
	Error() error
}
