package interfaces

import (
	"go-auth/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByPhoneNumber(phoneNumber string) (*models.User, error)
	GetUsers(page, limit int, search string) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}
