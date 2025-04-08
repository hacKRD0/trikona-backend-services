package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hacKRD0/trikona_go/user-management-service/config"
	"github.com/hacKRD0/trikona_go/user-management-service/database"
	"github.com/hacKRD0/trikona_go/user-management-service/handlers"
)

func main() {
	// Load config variables from config package
	cfg := config.LoadConfig()

	// Connect to the database
	database.ConnectDatabase()

	// Create a new Gin router
	router := gin.Default()

	// Configure CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FRONTEND_BASE_URL}, // your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Register routes
	router.GET("/", func (c * gin.Context) {
		c.JSON(200, gin.H{"message": "User Management Service is running!"})
	})

	// Register user routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/forgot-password", handlers.ForgotPassword)
		auth.POST("/reset-password", handlers.ResetPassword)
		auth.PUT("/update-user/:id", handlers.UpdateUser)
	}

	port := cfg.PORT

	// Start the server
	log.Printf("🚀 Starting User Management Service on port %s...", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}