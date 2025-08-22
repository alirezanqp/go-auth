package services

import (
	"fmt"
	"math"
	"time"

	"go-auth/internal/interfaces"
	"go-auth/internal/models"
	"go-auth/pkg/utils"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *UserService) GetUsers(page, limit int, search string) (*models.UsersListResponse, error) {
	if validationErrors := utils.ValidatePaginationParams(page, limit); validationErrors.HasErrors() {
		return nil, fmt.Errorf("validation failed: %s", validationErrors.Error())
	}

	if search != "" && !utils.IsValidSearchQuery(search) {
		return nil, fmt.Errorf("invalid search query")
	}

	users, total, err := s.userRepo.GetUsers(page, limit, search)
	if err != nil {
		return nil, err
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = models.UserResponse{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.UsersListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetUserStats() (map[string]interface{}, error) {
	users, _, err := s.userRepo.GetUsers(1, 1000000, "")
	if err != nil {
		return nil, err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	thisWeek := today.AddDate(0, 0, -7)
	thisMonth := today.AddDate(0, -1, 0)

	var todayCount, weekCount, monthCount int
	for _, user := range users {
		if user.CreatedAt.After(today) {
			todayCount++
		}
		if user.CreatedAt.After(thisWeek) {
			weekCount++
		}
		if user.CreatedAt.After(thisMonth) {
			monthCount++
		}
	}

	return map[string]interface{}{
		"total_users":      len(users),
		"users_today":      todayCount,
		"users_this_week":  weekCount,
		"users_this_month": monthCount,
		"timestamp":        now.Format(time.RFC3339),
	}, nil
}
