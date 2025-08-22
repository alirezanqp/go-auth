package handlers

import (
	"net/http"

	"go-auth/internal/config"
	"go-auth/internal/models"
	"go-auth/internal/services"
	"go-auth/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	otpService *services.OTPService
	config     *config.Config
}

func NewAuthHandler(otpService *services.OTPService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		otpService: otpService,
		config:     config,
	}
}

// @Summary Send OTP
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body models.SendOTPRequest true "Phone number"
// @Success 200 {object} models.SendOTPResponse
// @Router /auth/send-otp [post]
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req models.SendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := h.otpService.SendOTP(req.PhoneNumber); err != nil {
		appErr := utils.HandleError(err)
		c.JSON(appErr.HTTPCode, models.ErrorResponse{
			Success: false,
			Message: "Failed to send OTP",
			Error:   appErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, models.SendOTPResponse{
		Success: true,
		Message: "OTP sent successfully",
	})
}

// @Summary Verify OTP
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "Phone number and OTP"
// @Success 200 {object} models.VerifyOTPResponse
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.otpService.VerifyOTP(req.PhoneNumber, req.Code)
	if err != nil {
		appErr := utils.HandleError(err)
		c.JSON(appErr.HTTPCode, models.ErrorResponse{
			Success: false,
			Message: "OTP verification failed",
			Error:   appErr.Message,
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.PhoneNumber, h.config.JWT.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.VerifyOTPResponse{
		Success: true,
		Message: "Authentication successful",
		Token:   token,
		User:    user,
	})
}

// @Summary Get user profile
// @Tags authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	phoneNumber, _ := c.Get("phone_number")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user": gin.H{
			"id":           userID,
			"phone_number": phoneNumber,
		},
	})
}
