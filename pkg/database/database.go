package database

import (
	"fmt"
	"log"
	"os"

	"github.com/hacKRD0/trikona_go/internal/user-management-service/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
	// Construct the DSN
	env := os.Getenv("ENV")
	sslmode := "disable"
	if env == "prod" {
		sslmode = "require"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		sslmode,
	)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto migrate all the database schemas
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate database: %v", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
} 