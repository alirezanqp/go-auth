package middleware

import (
	"time"

	"go-auth/pkg/utils"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware creates a middleware for logging HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get request information
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Log the request
		utils.LogRequest(method, path, userAgent, clientIP, statusCode, latency)
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a simple request ID (in production, use UUID)
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + utils.GenerateRandomString(6)
}

// RecoveryWithLogging creates a recovery middleware that logs panics
func RecoveryWithLogging() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		utils.LogWithFields(map[string]interface{}{
			"panic":     recovered,
			"path":      c.Request.URL.Path,
			"method":    c.Request.Method,
			"client_ip": c.ClientIP(),
			"type":      "panic_recovery",
		}).Error("Panic recovered")

		c.AbortWithStatus(500)
	})
}
