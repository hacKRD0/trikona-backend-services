package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hacKRD0/trikona_go/user-management-service/config"
	"github.com/hacKRD0/trikona_go/user-management-service/database"
)

func main() {
	// Load config variables from config package
	cfg := config.LoadConfig()

	// Connect to the database
	database.ConnectDatabase()

	// Create a new Gin router
	router := gin.Default()

	// Register routes
	router.GET("/", func (c * gin.Context) {
		c.JSON(200, gin.H{"message": "User Management Service is running!"})
	})

	port := cfg.PORT

	// Start the server
	log.Printf("🚀 Starting User Management Service on port %s...", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}