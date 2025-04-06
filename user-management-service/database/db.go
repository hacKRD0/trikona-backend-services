package database

import (
	"fmt"
	"log"

	"github.com/hacKRD0/trikona_go/user-management-service/config"
	"github.com/hacKRD0/trikona_go/user-management-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load config variables from config package
	cfg := config.LoadConfig()

	// Construct the DSN (Data Source Name) for the database connection
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME)

	// Open a connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate the database:", err)
	}

	DB = db
	log.Println("✓✓✓ Connected to the database successfully!")
}