package integration_test_config

import (
	"log"
	"time"

	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	GormDB *gorm.DB
}

func GetTestDatabase(connStr string) *Database {
	dbInstance := &Database{}
	dbInstance.InitialiseTestDB(connStr)
	return dbInstance
}

func (db *Database) InitialiseTestDB(connStr string) {
	log.Printf("Connecting to Database with GORM...")

	var gormDB *gorm.DB

	for i := range 20 {
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

func (db *Database) CheckTestDatabaseHealth() error {
	sqlDB, err := db.GormDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
func (db *Database) CloseTestDatabase() error {
	sqlDB, err := db.GormDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *Database) RunTestMigrations() error {
	return db.GetTestGormDatabase().AutoMigrate(
		&userModel.User{},
	)
}

func (db *Database) GetTestGormDatabase() *gorm.DB {
	return db.GormDB
}
