package services

import (
	"fmt"
	"time"

	"go-auth/internal/config"
	"go-auth/internal/interfaces"
	"go-auth/internal/models"
	"go-auth/pkg/utils"
)

type OTPService struct {
	config         *config.Config
	otpRepo        interfaces.OTPRepository
	otpAttemptRepo interfaces.OTPAttemptRepository
	userRepo       interfaces.UserRepository
}

func NewOTPService(config *config.Config, otpRepo interfaces.OTPRepository, otpAttemptRepo interfaces.OTPAttemptRepository, userRepo interfaces.UserRepository) *OTPService {
	return &OTPService{
		config:         config,
		otpRepo:        otpRepo,
		otpAttemptRepo: otpAttemptRepo,
		userRepo:       userRepo,
	}
}

func (s *OTPService) SendOTP(phoneNumber string) error {
	if validationErrors := utils.ValidatePhoneNumber(phoneNumber); validationErrors.HasErrors() {
		utils.LogSecurityEvent("invalid_phone_number", "", phoneNumber, validationErrors.Error())
		return fmt.Errorf("validation failed: %s", validationErrors.Error())
	}

	if err := s.checkRateLimit(phoneNumber); err != nil {
		return err
	}

	otpCode, err := utils.GenerateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	expiresAt := time.Now().Add(s.config.OTP.ExpiryTime)
	otp := &models.OTP{
		PhoneNumber: phoneNumber,
		Code:        otpCode,
		ExpiresAt:   expiresAt,
		IsUsed:      false,
	}

	if err := s.otpRepo.Create(otp); err != nil {
		return err
	}

	attempt := &models.OTPAttempt{
		PhoneNumber: phoneNumber,
	}
	s.otpAttemptRepo.Create(attempt)

	utils.LogOTPGenerated(phoneNumber, otpCode, expiresAt)
	return nil
}

func (s *OTPService) VerifyOTP(phoneNumber, code string) (*models.User, error) {
	if validationErrors := utils.ValidatePhoneNumber(phoneNumber); validationErrors.HasErrors() {
		return nil, fmt.Errorf("validation failed: %s", validationErrors.Error())
	}

	if validationErrors := utils.ValidateOTPCode(code); validationErrors.HasErrors() {
		return nil, fmt.Errorf("validation failed: %s", validationErrors.Error())
	}

	otp, err := s.otpRepo.GetValidOTP(phoneNumber, code)
	if err != nil {
		utils.LogOTPVerification(phoneNumber, code, false, err.Error())
		return nil, err
	}

	if err := s.otpRepo.MarkAsUsed(otp.ID); err != nil {
		return nil, err
	}

	utils.LogOTPVerification(phoneNumber, code, true, "OTP verified successfully")

	user, err := s.userRepo.GetByPhoneNumber(phoneNumber)
	if err != nil {
		if err == utils.ErrUserNotFound {
			newUser := &models.User{
				PhoneNumber: phoneNumber,
			}
			if err := s.userRepo.Create(newUser); err != nil {
				return nil, err
			}
			utils.LogUserRegistration(newUser.ID.String(), phoneNumber)
			return newUser, nil
		}
		return nil, err
	}

	utils.LogUserLogin(user.ID.String(), phoneNumber)
	return user, nil
}

func (s *OTPService) checkRateLimit(phoneNumber string) error {
	cutoffTime := time.Now().Add(-s.config.OTP.RateWindow)
	count, err := s.otpAttemptRepo.CountRecentAttempts(phoneNumber, cutoffTime)
	if err != nil {
		return err
	}

	if count >= int64(s.config.OTP.MaxAttempts) {
		utils.LogRateLimit(phoneNumber, int(count), s.config.OTP.MaxAttempts)
		return utils.ErrRateLimitExceeded
	}

	return nil
}

func (s *OTPService) CleanupExpiredOTPs() error {
	return s.otpRepo.DeleteExpired()
}
