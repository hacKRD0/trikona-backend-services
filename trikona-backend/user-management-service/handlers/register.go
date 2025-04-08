package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hacKRD0/trikona_go/user-management-service/database"
	"github.com/hacKRD0/trikona_go/user-management-service/models"
	"github.com/mailjet/mailjet-apiv3-go"
	"golang.org/x/crypto/bcrypt"
)

func sendEmailVerification(toEmail string, token string) error {
	registerURL := fmt.Sprintf("%s/register?token=%s", cfg.FRONTEND_BASE_URL, token)

	log.Printf("Sending password reset email to %s from %s", toEmail, cfg.EMAIL_FROM)
	log.Printf("API key: %s\nAPI secret: %s", cfg.MAILJET_API_KEY, cfg.MAILJET_API_SECRET)
	// Initialize the Mailjet client with your API key and secret.
	mjClient := mailjet.NewMailjetClient(cfg.MAILJET_API_KEY, cfg.MAILJET_API_SECRET)

	// Create a new email message.
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: cfg.EMAIL_FROM,
				Name:  "Trikona Support",
			},
			To: &mailjet.RecipientsV31{
          mailjet.RecipientV31 {
					Email: toEmail,
				},
			},
			Subject:  "Verify Your Email",
			TextPart: fmt.Sprintf("Hello,\n\nPlease verify your email address by clicking the link below:\n\n%s\n\nThanks,\nTrikona Support Team", registerURL),
			HTMLPart: fmt.Sprintf("<p>Hello,</p><p>Please click the link below to verify your email address:</p><p><a href=\"%s\">Verify Email</a></p><p>Thanks,<br>Trikona Support Team</p>", registerURL),
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}

	// Send the email using Mailjet's API.
	response, err := mjClient.SendMailV31(&messages)
	if err != nil {
		fmt.Printf("Failed to send reset email: %v\n", err)
		return err
	}

	fmt.Printf("Password reset email sent. Response data: %v\n", response)
	return nil

}

func RequestVerification(c *gin.Context) {
	// Input DTO struct for user signup
	var input struct {
		Email string `form:"email" binding:"required,email"`
	}

	// Bind JSON input
	if err := c.ShouldBind(&input); err != nil {
		log.Printf("Error binding JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists in the database
	if database.DB.Where("email = ?", input.Email).First(&models.User{}).RowsAffected > 0 {
		log.Printf("Email %s already exists", input.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists!"})
		return
	}

	// Send verification email with token
	verificationToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": input.Email,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := verificationToken.SignedString([]byte(os.Getenv("JWT_REGISTRATION_SECRET")))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send verification email
	if err := sendEmailVerification(input.Email, tokenString); err != nil {
		log.Printf("Error sending verification email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

func Register(c *gin.Context) {
	// Input DTO struct for user signup
	var input struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		Role      string `json:"role"` 
		Token     string `json:"token" binding:"required"`
	} 

	// Bind JSON input
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		log.Printf("Error binding JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the token
	token, err := jwt.Parse(input.Token, func (t* jwt.Token) (interface{}, error) {
		// Ensure that the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		// Return the secret key
		return []byte(os.Getenv("JWT_REGISTRATION_SECRET")), nil
	})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token!"})
		return
	}

	// Check if the token is valid and if the email matches
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		log.Println("Invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token!"})
		return
	}

	// Get the email from the token
	email := token.Claims.(jwt.MapClaims)["email"].(string)
	if email != input.Email {
		log.Printf("Email in token does not match email in input: %s != %s", email, input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token!"})
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

	// Generate new JWT token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"role": user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the JWT token as a cookie
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfuly!", "token": tokenString})
}