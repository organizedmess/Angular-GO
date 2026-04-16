package db

import (
	"log"
	"os"

	"url-shortener/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=127.0.0.1 user=postgres password=postgres dbname=url_shortener port=5434 sslmode=disable"
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := database.AutoMigrate(&models.URL{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return database
}
