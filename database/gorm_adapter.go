package database

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type GormAdapter struct {
	db *gorm.DB
}

func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{db: db}
}

func (g *GormAdapter) Where(query any, args ...any) IDatabase {
	return &GormAdapter{db: g.db.Where(query, args...)}
}

func (g *GormAdapter) First(dest any, conds ...any) IDatabase {
	return &GormAdapter{db: g.db.First(dest, conds...)}
}

func (g *GormAdapter) Find(dest any, conds ...any) IDatabase {
	return &GormAdapter{db: g.db.Find(dest, conds...)}
}

func (g *GormAdapter) Select(query any, args ...any) IDatabase {
	return &GormAdapter{db: g.db.Select(query, args...)}
}

func (g *GormAdapter) Create(value any) IDatabase {
	return &GormAdapter{db: g.db.Create(value)}
}

func (g *GormAdapter) Delete(value any, conds ...any) IDatabase {
	return &GormAdapter{db: g.db.Delete(value, conds...)}
}

func (g *GormAdapter) Model(value any) IDatabase {
	return &GormAdapter{db: g.db.Model(value)}
}

func (g *GormAdapter) Updates(values any) IDatabase {
	return &GormAdapter{db: g.db.Updates(values)}
}

func (g *GormAdapter) Limit(limit int) IDatabase {
	return &GormAdapter{db: g.db.Limit(limit)}
}

func (g *GormAdapter) Offset(offset int) IDatabase {
	return &GormAdapter{db: g.db.Offset(offset)}
}

func (g *GormAdapter) Order(value any) IDatabase {
	return &GormAdapter{db: g.db.Order(value)}
}

func (g *GormAdapter) Count(count *int64) IDatabase {
	return &GormAdapter{db: g.db.Count(count)}
}

func (g *GormAdapter) WithContext(ctx context.Context) IDatabase {
	return &GormAdapter{db: g.db.WithContext(ctx)}
}

func (g *GormAdapter) Begin(opts ...*sql.TxOptions) IDatabase {
	return &GormAdapter{db: g.db.Begin(opts...)}
}

func (g *GormAdapter) Commit() IDatabase {
	return &GormAdapter{db: g.db.Commit()}
}

func (g *GormAdapter) Rollback() IDatabase {
	return &GormAdapter{db: g.db.Rollback()}
}

func (g *GormAdapter) DB() (*sql.DB, error) {
	return g.db.DB()
}

func (g *GormAdapter) AutoMigrate(dst ...any) error {
	return g.db.AutoMigrate(dst...)
}

// Error implements GORM's error access pattern
func (g *GormAdapter) Error() error {
	return g.db.Error
}
