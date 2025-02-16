package database

import (
	"fmt"
	"log"

	"github.com/hr0mk4/test_api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := "host=db user=postgres password=secret dbname=merch_store port=5432 sslmode=disable TimeZone=UTC"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return err
	}

	if err := DB.AutoMigrate(&models.User{}, &models.Purchase{}, &models.Transaction{}); err != nil {
		log.Printf("Failed to migrate: %v", err)
		return err
	}

	fmt.Println("Database connected")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
