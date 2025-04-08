package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mailjet/mailjet-apiv3-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/hacKRD0/trikona_go/user-management-service/config"
	"github.com/hacKRD0/trikona_go/user-management-service/database"
	"github.com/hacKRD0/trikona_go/user-management-service/models"
)

var cfg = config.LoadConfig()

func sendPasswordResetEmail(toEmail string, firstName string, token string) error {
	// Construct the password reset URL.
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", cfg.FRONTEND_BASE_URL, token)

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
			Subject:  "Password Reset Request",
			TextPart: fmt.Sprintf("Hello %s,\n\nWe received a request to reset your password. Please use the link below to reset your password:\n\n%s\n\nThe link will be valid for the next 3 hours. If you did not request this, please ignore this email.\n\nThanks,\nTrikona Support Team", firstName, resetURL),
			HTMLPart: fmt.Sprintf("<p>Hello,</p><p>We received a request to reset your password. Please click the link below to reset your password:</p><p><a href=\"%s\">Reset Password</a></p><p>The link will be valid for the next 3 hours. If you did not request this, please ignore this email.</p><p>Thanks,<br>Trikona Support Team</p>", resetURL),
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

func ForgotPassword(c *gin.Context) {
	// Get the email from the request JSON body
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBind(&input); err != nil {
		log.Printf("Error binding JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Check if email exists in the database
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {	
		log.Printf("Email %s does not exist", input.Email)
		c.JSON(http.StatusNotFound, gin.H{"error": "Email does not exist!"})
		return
	}

	// Send password reset email with token
	resetToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := resetToken.SignedString([]byte(cfg.JWT_RESET_SECRET))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	// TODO: Send password reset email with token
	if err := sendPasswordResetEmail(user.Email, user.FirstName, tokenString); err != nil {
		log.Printf("Error sending password reset email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent!"})
}

func ResetPassword(c *gin.Context) {
	// Get the token from the request JSON body
	var input struct {
		Token string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBind(&input); err != nil {
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
		return []byte(cfg.JWT_RESET_SECRET), nil
	})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token!"})
		return
	}

	// Check if the token claims are valid
	claims, ok := token.Claims.(jwt.MapClaims) 
	if !ok {
		log.Printf("Error parsing token claims: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims!"})
		return
	}

	// Get the user ID from the token claims
	userID := claims["user_id"].(float64)

	// Check if the user exists in the database
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("User with ID %f does not exist", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist!"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the user password in the database
	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		log.Printf("Error updating user password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful!"})
}