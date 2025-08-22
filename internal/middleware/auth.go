package middleware

import (
	"net/http"
	"strings"

	"go-auth/internal/config"
	"go-auth/internal/models"
	"go-auth/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Authorization header required",
				Error:   "missing_auth_header",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Invalid authorization format. Use: Bearer <token>",
				Error:   "invalid_auth_format",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Token is required",
				Error:   "missing_token",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(token, cfg.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Invalid or expired token",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("phone_number", claims.PhoneNumber)
		c.Set("claims", claims)

		c.Next()
	}
}

func OptionalAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")

			if claims, err := utils.ValidateJWT(token, cfg.JWT.Secret); err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("phone_number", claims.PhoneNumber)
				c.Set("claims", claims)
				c.Set("authenticated", true)
			}
		}

		c.Next()
	}
}
 