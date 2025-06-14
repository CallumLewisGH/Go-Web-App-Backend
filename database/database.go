package database

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	GormDB *gorm.DB
}

var (
	dbInstance  *Database
	once        sync.Once
	testMode    bool
	testConnStr string
)

// SetTestMode enables test mode with a specific connection string
func SetTestMode(connStr string) {
	testMode = true
	testConnStr = connStr
	dbInstance = nil
	once = sync.Once{}
}

// SetProdMode switches back to production mode
func SetProdMode() {
	testMode = false
	testConnStr = ""
	dbInstance = nil
	once = sync.Once{}
}

func GetDatabase() *Database {
	once.Do(func() {
		dbInstance = &Database{}
		if testMode {
			dbInstance.InitialiseTestDB(testConnStr)
		} else {
			dbInstance.InitialiseDB()
		}
	})
	return dbInstance
}

func (db *Database) InitialiseTestDB(connStr string) {
	log.Printf("Connecting to Database with GORM...")

	var gormDB *gorm.DB

	for i := range 5 {
		if i > 0 {
			time.Sleep(1 * time.Second)
			log.Printf("Retrying connection to database, attempt %d", i+1)
		}

		gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
			PrepareStmt: true,
		})
		if err == nil {
			db.GormDB = gormDB
			log.Printf("GORM Database Connection Succeeded")
			return
		}
		log.Printf("Failed to open database: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	db.GormDB = gormDB
	log.Printf("GORM Database Connection Succeeded")
}

func (db *Database) InitialiseDB() {
	log.Printf("Connecting to Database with GORM...")
	err := godotenv.Load("/home/callum/Desktop/Go-Web-App-Backend/.dev.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	var gormDB *gorm.DB

	connStr := os.Getenv("DATABASE_CONNECTION_STRING")

	for i := range 5 {
		if i > 0 {
			time.Sleep(1 * time.Second)
			log.Printf("Retrying connection to database, attempt %d", i+1)
		}

		gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
			PrepareStmt: true,
		})
		if err == nil {
			db.GormDB = gormDB
			log.Printf("GORM Database Connection Succeeded")
			return
		}
		log.Printf("Failed to open database: %v", err)
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

func (db *Database) GetGormDatabase() *gorm.DB {
	return db.GormDB
}

func (db *Database) RunMigrations() error {
	for _, model := range GetModelRegistry().models {
		if err := GetDatabase().GetGormDatabase().AutoMigrate(model); err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) ClearAllTables() error {
	for _, model := range GetModelRegistry().models {
		if err := db.GormDB.Unscoped().Where("1 = 1").Delete(model).Error; err != nil {
			return err
		}
	}
	return nil
}
