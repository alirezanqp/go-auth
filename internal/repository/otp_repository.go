package repository

import (
	"errors"
	"fmt"
	"time"

	"go-auth/internal/interfaces"
	"go-auth/internal/models"
	"go-auth/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) interfaces.OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(otp *models.OTP) error {
	if err := r.db.Create(otp).Error; err != nil {
		utils.LogDatabaseOperation("create", "otps", false, err.Error())
		return fmt.Errorf("failed to create OTP: %w", err)
	}

	utils.LogDatabaseOperation("create", "otps", true, "")
	return nil
}

func (r *otpRepository) GetValidOTP(phoneNumber, code string) (*models.OTP, error) {
	var otp models.OTP
	err := r.db.Where("phone_number = ? AND code = ? AND is_used = false", phoneNumber, code).
		First(&otp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrInvalidOTP
		}
		utils.LogDatabaseOperation("find", "otps", false, err.Error())
		return nil, fmt.Errorf("failed to get OTP: %w", err)
	}

	if otp.IsExpired() {
		return nil, utils.ErrOTPExpired
	}

	return &otp, nil
}

func (r *otpRepository) MarkAsUsed(id uuid.UUID) error {
	result := r.db.Model(&models.OTP{}).Where("id = ?", id).Update("is_used", true)

	if result.Error != nil {
		utils.LogDatabaseOperation("update", "otps", false, result.Error.Error())
		return fmt.Errorf("failed to mark OTP as used: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return utils.ErrInvalidOTP
	}

	utils.LogDatabaseOperation("update", "otps", true, "")
	return nil
}

func (r *otpRepository) DeleteExpired() error {
	result := r.db.Where("expires_at < ?", time.Now()).Delete(&models.OTP{})

	if result.Error != nil {
		utils.LogDatabaseOperation("cleanup", "otps", false, result.Error.Error())
		return fmt.Errorf("failed to cleanup expired OTPs: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		utils.LogWithFields(map[string]interface{}{
			"rows_affected": result.RowsAffected,
			"type":          "cleanup",
			"table":         "otps",
		}).Info("Cleaned up expired OTPs")
	}

	return nil
}

type otpAttemptRepository struct {
	db *gorm.DB
}

func NewOTPAttemptRepository(db *gorm.DB) interfaces.OTPAttemptRepository {
	return &otpAttemptRepository{db: db}
}

func (r *otpAttemptRepository) Create(attempt *models.OTPAttempt) error {
	if err := r.db.Create(attempt).Error; err != nil {
		utils.LogDatabaseOperation("create", "otp_attempts", false, err.Error())
		return fmt.Errorf("failed to create OTP attempt: %w", err)
	}

	return nil
}

func (r *otpAttemptRepository) CountRecentAttempts(phoneNumber string, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.OTPAttempt{}).
		Where("phone_number = ? AND attempt_time > ?", phoneNumber, since).
		Count(&count).Error

	if err != nil {
		utils.LogDatabaseOperation("count", "otp_attempts", false, err.Error())
		return 0, fmt.Errorf("failed to count OTP attempts: %w", err)
	}

	return count, nil
}

func (r *otpAttemptRepository) DeleteOldAttempts(before time.Time) error {
	result := r.db.Where("attempt_time < ?", before).Delete(&models.OTPAttempt{})

	if result.Error != nil {
		utils.LogDatabaseOperation("cleanup", "otp_attempts", false, result.Error.Error())
		return fmt.Errorf("failed to cleanup old OTP attempts: %w", result.Error)
	}

	return nil
}
