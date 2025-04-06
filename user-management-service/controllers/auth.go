package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hacKRD0/trikona_go/user-management-service/database"
	"github.com/hacKRD0/trikona_go/user-management-service/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// Input DTO struct for user signup
	var input struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		Role      string `json:"role"` 
	} 

	// Bind JSON input
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		log.Printf("Error binding JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a new user in the database
	user := models.User{
		FirstName : input.FirstName,
		LastName : input.LastName,
		Email : input.Email,
		Password : string(hashedPassword),
		Role : input.Role, // Set the role to "guest" by default
	}

	// Check if the email already exists in the database
	if database.DB.Where("email = ?", input.Email).First(&user).RowsAffected > 0 {
		log.Printf("Email %s already exists", input.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Save the user to the database
	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"role": user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the JWT token as a cookie
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfuly!", "token": tokenString})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON input
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		log.Printf("Error binding JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the email and password from the request JSON body
	Email := c.PostForm("Email")
	Password := c.PostForm("Password")

	var user models.User
	// Check if the Email exists in the database
	if err := database.DB.Where("email = ?", Email).First(&user).Error; err != nil {
		log.Printf("Email %s does not exist", Email)
		c.JSON(http.StatusConflict, gin.H{"error": "Email does not exist!"})
		return
	}

	// Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(user.Password)); err != nil {
		log.Println("Password is incorrect")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect!"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the JWT token as a cookie
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": tokenString})
}