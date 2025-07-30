package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Alfian57/belajar-golang/internal/config"
	"github.com/Alfian57/belajar-golang/internal/logger"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(config config.DatabaseConfig) {
	// Initialize the database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.Host, config.Username, config.Password, config.Name, config.Port)

	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open a new database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logger.Log.Fatalf("error opening database connection: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Log.Fatalf("error opening database connection: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Log.Infoln("successfully connected to the database")
}
