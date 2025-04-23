package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hacKRD0/trikona_go/pkg/logger"
	"go.uber.org/zap"
)

// RequestLogger returns a middleware that logs request details
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID if not present
		requestID := c.GetString("request_id")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("request_id", requestID)
		}

		// Capture request details
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Log request start
		logger.Info("Request started",
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
		)

		// Process request
		c.Next()

		// Log request completion
		latency := time.Since(start)
		status := c.Writer.Status()
		responseSize := c.Writer.Size()

		logger.Info("Request completed",
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Int("status", status),
			zap.Int("response_size", responseSize),
			zap.Duration("latency", latency),
		)
	}
} 