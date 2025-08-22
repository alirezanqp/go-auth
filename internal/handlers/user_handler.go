package handlers

import (
	"net/http"
	"strconv"
	"time"

	"go-auth/internal/models"
	"go-auth/internal/services"
	"go-auth/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Summary Get user by ID
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid user ID format",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		appErr := utils.HandleError(err)
		c.JSON(appErr.HTTPCode, models.ErrorResponse{
			Success: false,
			Message: "Failed to get user",
			Error:   appErr.Message,
		})
		return
	}

	userResponse := models.UserResponse{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    userResponse,
	})
}

// @Summary Get users list
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param search query string false "Search by phone number"
// @Success 200 {object} models.UsersListResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	search := c.Query("search")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	response, err := h.userService.GetUsers(page, limit, search)
	if err != nil {
		appErr := utils.HandleError(err)
		c.JSON(appErr.HTTPCode, models.ErrorResponse{
			Success: false,
			Message: "Failed to get users",
			Error:   appErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// @Summary Get user statistics
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /users/stats [get]
func (h *UserHandler) GetUserStats(c *gin.Context) {
	stats, err := h.userService.GetUserStats()
	if err != nil {
		appErr := utils.HandleError(err)
		c.JSON(appErr.HTTPCode, models.ErrorResponse{
			Success: false,
			Message: "Failed to get user statistics",
			Error:   appErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats":   stats,
	})
}
