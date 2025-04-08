package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/hacKRD0/trikona_go/user-management-service/database"
	"github.com/hacKRD0/trikona_go/user-management-service/models"
)

// LinkedInTokenResponse represents the JSON response from LinkedIn token exchange.
type LinkedInTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// LinkedInProfile represents the JSON response from LinkedIn's profile endpoint.
type LinkedInProfile struct {
	ID        string `json:"id"`
	LocalizedFirstName string `json:"localizedFirstName"`
	LocalizedLastName  string `json:"localizedLastName"`
}

// LinkedInEmailResponse represents the JSON response from LinkedIn's email endpoint.
type LinkedInEmailResponse struct {
	Elements []struct {
		Handle struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"handle~"`
	} `json:"elements"`
}

// LinkedInLoginHandler handles LinkedIn login/signup.
func LinkedInLoginHandler(c *gin.Context) {
	var input struct {
		Code  string `json:"code" binding:"required"`
		State string `json:"state"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Exchange authorization code for an access token.
	tokenURL := "https://www.linkedin.com/oauth/v2/accessToken"
	data := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		input.Code,
		cfg.LINKEDIN_REDIRECT_URI,
		cfg.LINKEDIN_CLIENT_ID,
		cfg.LINKEDIN_CLIENT_SECRET,
	)
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token request"})
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to request access token"})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read token response"})
		return
	}

	var tokenRes LinkedInTokenResponse
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
		return
	}

	accessToken := tokenRes.AccessToken

	// Get user profile from LinkedIn.
	profileReq, err := http.NewRequest("GET", "https://api.linkedin.com/v2/me", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile request"})
		return
	}
	profileReq.Header.Set("Authorization", "Bearer "+accessToken)
	profileResp, err := client.Do(profileReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}
	defer profileResp.Body.Close()
	profileBody, err := io.ReadAll(profileResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read profile response"})
		return
	}

	var profile LinkedInProfile
	if err := json.Unmarshal(profileBody, &profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse profile response"})
		return
	}

	// Get user email from LinkedIn.
	emailReq, err := http.NewRequest("GET", "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email request"})
		return
	}
	emailReq.Header.Set("Authorization", "Bearer "+accessToken)
	emailResp, err := client.Do(emailReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get email"})
		return
	}
	defer emailResp.Body.Close()
	emailBody, err := io.ReadAll(emailResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read email response"})
		return
	}
	var emailRes LinkedInEmailResponse
	if err := json.Unmarshal(emailBody, &emailRes); err != nil || len(emailRes.Elements) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse email response"})
		return
	}
	emailAddress := emailRes.Elements[0].Handle.EmailAddress

	// Check if the user exists in the database; if not, create a new user.
	var user models.User
	err = database.DB.Where("email = ?", emailAddress).First(&user).Error
	if err != nil {
		// Assuming record not found; create a new user.
		user = models.User{
			FirstName: profile.LocalizedFirstName,
			LastName:  profile.LocalizedLastName,
			Email:     emailAddress,
			// For social logins, you might set a default or random password,
			// or use a different field to indicate social login.
			Password:  "", // Alternatively, store a placeholder.
			Role:      models.RoleGuest, // Or another default role.
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Generate a JWT for the user (for example, for session management).
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	jwtToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "LinkedIn login successful",
		"token":   jwtToken,
	})
}
