package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CorsConfig holds the CORS configuration
type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCorsConfig returns a default CORS configuration
func DefaultCorsConfig() *CorsConfig {
	return &CorsConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}
}

// Cors returns a CORS middleware handler
func Cors(config *CorsConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCorsConfig()
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		// Check if the origin is allowed
		isAllowed := false
		for _, allowedOrigin := range config.AllowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))
		c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ","))
		c.Writer.Header().Set("Access-Control-Max-Age", string(config.MaxAge))

		if config.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
} 