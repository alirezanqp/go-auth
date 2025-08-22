package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}

	var messages []string
	for _, err := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}

	return strings.Join(messages, ", ")
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

func ValidatePhoneNumber(phoneNumber string) ValidationErrors {
	var errors ValidationErrors

	if phoneNumber == "" {
		errors = append(errors, ValidationError{
			Field:   "phone_number",
			Message: "phone number is required",
		})
		return errors
	}

	phoneNumber = strings.TrimSpace(phoneNumber)

	if len(phoneNumber) < 10 {
		errors = append(errors, ValidationError{
			Field:   "phone_number",
			Message: "phone number must be at least 10 digits",
		})
	}

	if len(phoneNumber) > 15 {
		errors = append(errors, ValidationError{
			Field:   "phone_number",
			Message: "phone number must not exceed 15 digits",
		})
	}

	if !phoneRegex.MatchString(phoneNumber) {
		errors = append(errors, ValidationError{
			Field:   "phone_number",
			Message: "invalid phone number format",
		})
	}

	return errors
}

func ValidateOTPCode(code string) ValidationErrors {
	var errors ValidationErrors

	if code == "" {
		errors = append(errors, ValidationError{
			Field:   "code",
			Message: "OTP code is required",
		})
		return errors
	}

	code = strings.TrimSpace(code)

	if len(code) != 6 {
		errors = append(errors, ValidationError{
			Field:   "code",
			Message: "OTP code must be exactly 6 digits",
		})
	}

	for _, char := range code {
		if !unicode.IsDigit(char) {
			errors = append(errors, ValidationError{
				Field:   "code",
				Message: "OTP code must contain only digits",
			})
			break
		}
	}

	return errors
}

func ValidateUserID(userID string) ValidationErrors {
	var errors ValidationErrors

	if userID == "" {
		errors = append(errors, ValidationError{
			Field:   "user_id",
			Message: "user ID is required",
		})
		return errors
	}

	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(userID) {
		errors = append(errors, ValidationError{
			Field:   "user_id",
			Message: "invalid user ID format",
		})
	}

	return errors
}

func ValidatePaginationParams(page, limit int) ValidationErrors {
	var errors ValidationErrors

	if page < 1 {
		errors = append(errors, ValidationError{
			Field:   "page",
			Message: "page must be greater than 0",
		})
	}

	if limit < 1 {
		errors = append(errors, ValidationError{
			Field:   "limit",
			Message: "limit must be greater than 0",
		})
	}

	if limit > 100 {
		errors = append(errors, ValidationError{
			Field:   "limit",
			Message: "limit cannot exceed 100",
		})
	}

	return errors
}

func SanitizeString(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, "\t", "")

	return input
}

func IsValidSearchQuery(query string) bool {
	if len(query) < 2 {
		return false
	}

	if len(query) > 50 {
		return false
	}

	maliciousPatterns := []string{
		"<script",
		"javascript:",
		"onload=",
		"onerror=",
		"DROP TABLE",
		"DELETE FROM",
		"INSERT INTO",
		"UPDATE SET",
	}

	lowerQuery := strings.ToLower(query)
	for _, pattern := range maliciousPatterns {
		if strings.Contains(lowerQuery, strings.ToLower(pattern)) {
			return false
		}
	}

	return true
}
