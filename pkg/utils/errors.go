package utils

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	HTTPCode int    `json:"-"`
	Details  string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	ErrInvalidPhoneNumber = &AppError{
		Code:     "INVALID_PHONE_NUMBER",
		Message:  "Invalid phone number format",
		HTTPCode: http.StatusBadRequest,
	}

	ErrInvalidOTP = &AppError{
		Code:     "INVALID_OTP",
		Message:  "Invalid or expired OTP",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrOTPExpired = &AppError{
		Code:     "OTP_EXPIRED",
		Message:  "OTP has expired",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrRateLimitExceeded = &AppError{
		Code:     "RATE_LIMIT_EXCEEDED",
		Message:  "Too many requests. Please try again later",
		HTTPCode: http.StatusTooManyRequests,
	}

	ErrUserNotFound = &AppError{
		Code:     "USER_NOT_FOUND",
		Message:  "User not found",
		HTTPCode: http.StatusNotFound,
	}

	ErrUnauthorized = &AppError{
		Code:     "UNAUTHORIZED",
		Message:  "Authentication required",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrInvalidToken = &AppError{
		Code:     "INVALID_TOKEN",
		Message:  "Invalid or expired token",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrInternalServer = &AppError{
		Code:     "INTERNAL_SERVER_ERROR",
		Message:  "Internal server error",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrInvalidRequest = &AppError{
		Code:     "INVALID_REQUEST",
		Message:  "Invalid request format",
		HTTPCode: http.StatusBadRequest,
	}

	ErrValidationFailed = &AppError{
		Code:     "VALIDATION_FAILED",
		Message:  "Validation failed",
		HTTPCode: http.StatusBadRequest,
	}
)

func NewAppError(code, message string, httpCode int) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
	}
}

func (e *AppError) WithDetails(details string) *AppError {
	return &AppError{
		Code:     e.Code,
		Message:  e.Message,
		HTTPCode: e.HTTPCode,
		Details:  details,
	}
}

func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

func HandleError(err error) *AppError {
	if appErr, ok := IsAppError(err); ok {
		return appErr
	}

	return ErrInternalServer.WithDetails(err.Error())
}
