package interfaces

import (
	"time"

	"go-auth/internal/models"

	"github.com/google/uuid"
)

type OTPRepository interface {
	Create(otp *models.OTP) error
	GetValidOTP(phoneNumber, code string) (*models.OTP, error)
	MarkAsUsed(id uuid.UUID) error
	DeleteExpired() error
}

type OTPAttemptRepository interface {
	Create(attempt *models.OTPAttempt) error
	CountRecentAttempts(phoneNumber string, since time.Time) (int64, error)
	DeleteOldAttempts(before time.Time) error
}
