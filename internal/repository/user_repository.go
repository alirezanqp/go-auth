package repository

import (
	"fmt"

	"go-auth/internal/interfaces"
	"go-auth/internal/models"
	"go-auth/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		utils.LogDatabaseOperation("create", "users", false, err.Error())
		return fmt.Errorf("failed to create user: %w", err)
	}

	utils.LogDatabaseOperation("create", "users", true, "")
	return nil
}

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrUserNotFound
		}
		utils.LogDatabaseOperation("find", "users", false, err.Error())
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrUserNotFound
		}
		utils.LogDatabaseOperation("find", "users", false, err.Error())
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetUsers(page, limit int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	if search != "" {
		query = query.Where("phone_number ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		utils.LogDatabaseOperation("count", "users", false, err.Error())
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		utils.LogDatabaseOperation("find", "users", false, err.Error())
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

func (r *userRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		utils.LogDatabaseOperation("update", "users", false, err.Error())
		return fmt.Errorf("failed to update user: %w", err)
	}

	utils.LogDatabaseOperation("update", "users", true, "")
	return nil
}

func (r *userRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.User{}, id)

	if result.Error != nil {
		utils.LogDatabaseOperation("delete", "users", false, result.Error.Error())
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return utils.ErrUserNotFound
	}

	utils.LogDatabaseOperation("delete", "users", true, "")
	return nil
}
