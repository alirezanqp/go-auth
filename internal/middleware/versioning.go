package middleware

import (
	"net/http"
	"strings"

	"go-auth/internal/models"
	"go-auth/pkg/utils"

	"github.com/gin-gonic/gin"
)

// APIVersionMiddleware handles API versioning through headers or URL
func APIVersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get version from header first
		version := c.GetHeader("API-Version")

		// If not in header, try to extract from URL path
		if version == "" {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api/v") {
				parts := strings.Split(path, "/")
				if len(parts) >= 3 {
					version = parts[2] // e.g., "v1" from "/api/v1/auth/login"
				}
			}
		}

		// Default to v1 if no version specified
		if version == "" {
			version = "v1"
		}

		// Validate version
		if !isValidVersion(version) {
			utils.LogSecurityEvent("invalid_api_version", "", "", "Invalid API version requested: "+version)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Message: "Invalid API version",
				Error:   "Supported versions: v1",
			})
			c.Abort()
			return
		}

		// Set version in context for handlers to use
		c.Set("api_version", version)

		// Add version to response headers
		c.Header("API-Version", version)

		c.Next()
	}
}

// DeprecationWarningMiddleware adds deprecation warnings for old API versions
func DeprecationWarningMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		version, exists := c.Get("api_version")
		if !exists {
			c.Next()
			return
		}

		versionStr := version.(string)

		// Check if version is deprecated
		if isDeprecatedVersion(versionStr) {
			c.Header("Deprecation", "true")
			c.Header("Sunset", "2024-12-31") // Example sunset date
			c.Header("Link", "</api/v1>; rel=\"successor-version\"")

			utils.LogWithFields(map[string]interface{}{
				"deprecated_version": versionStr,
				"client_ip":          c.ClientIP(),
				"user_agent":         c.Request.UserAgent(),
				"path":               c.Request.URL.Path,
				"type":               "deprecated_api_usage",
			}).Warn("Deprecated API version used")
		}

		c.Next()
	}
}

// isValidVersion checks if the provided version is supported
func isValidVersion(version string) bool {
	supportedVersions := []string{"v1"}

	for _, v := range supportedVersions {
		if v == version {
			return true
		}
	}

	return false
}

// isDeprecatedVersion checks if the provided version is deprecated
func isDeprecatedVersion(version string) bool {
	deprecatedVersions := []string{
		// Add deprecated versions here, e.g., "v0"
	}

	for _, v := range deprecatedVersions {
		if v == version {
			return true
		}
	}

	return false
}
