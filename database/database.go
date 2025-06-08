package database

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/CallumLewisGH/Generic-Service-Base/internal/domain/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	GormDB *gorm.DB
}

var (
	dbInstance *Database
	once       sync.Once
)

func GetDatabase() *Database {
	once.Do(func() {
		dbInstance = &Database{}
		dbInstance.InitialiseDB()
	})
	return dbInstance
}

func (db *Database) InitialiseDB() {
	log.Printf("Connecting to Database with GORM...")
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DATABASE_CONNECTION_STRING")

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	db.GormDB = gormDB
	log.Printf("GORM Database Connection Succeeded")
}

func (db *Database) CheckDatabaseHealth() error {
	sqlDB, err := db.GormDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (db *Database) CloseDatabase() error {
	sqlDB, err := db.GormDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *Database) RunMigrations() error {
	return db.GetGormDatabase().AutoMigrate(&models.User{})
}

func (db *Database) GetGormDatabase() *gorm.DB {
	return db.GormDB
}
