package models

import (
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PhoneNumber string    `json:"phone_number" gorm:"index;not null"`
	Code        string    `json:"code" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	ExpiresAt   time.Time `json:"expires_at" gorm:"not null"`
	IsUsed      bool      `json:"is_used" gorm:"default:false"`
}

func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func (OTP) TableName() string {
	return "otps"
}

type OTPAttempt struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PhoneNumber string    `json:"phone_number" gorm:"index;not null"`
	AttemptTime time.Time `json:"attempt_time" gorm:"autoCreateTime"`
}

func (OTPAttempt) TableName() string {
	return "otp_attempts"
}
