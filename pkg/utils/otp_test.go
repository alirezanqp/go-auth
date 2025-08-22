package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateOTP(t *testing.T) {
	otp, err := GenerateOTP()

	assert.NoError(t, err)
	assert.Len(t, otp, 6)

	for _, char := range otp {
		assert.True(t, char >= '0' && char <= '9', "OTP should contain only digits")
	}
}

func TestGenerateRandomString(t *testing.T) {
	tests := []int{5, 10, 16, 32}

	for _, length := range tests {
		t.Run(fmt.Sprintf("length_%d", length), func(t *testing.T) {
			result := GenerateRandomString(length)
			assert.Len(t, result, length)
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		expectError bool
	}{
		{"Valid international", "+1234567890", false},
		{"Valid local", "1234567890", false},
		{"Too short", "123", true},
		{"Too long", "123456789012345678", true},
		{"Contains letters", "+123abc456", true},
		{"Contains special chars", "+123-456-789", true},
		{"Empty string", "", true},
		{"Only plus", "+", true},
		{"Valid Iranian", "+989123456789", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidatePhoneNumber(tt.phoneNumber)
			hasErrors := errors.HasErrors()

			if tt.expectError {
				assert.True(t, hasErrors, "Expected validation error for %s", tt.phoneNumber)
			} else {
				assert.False(t, hasErrors, "Expected no validation error for %s", tt.phoneNumber)
			}
		})
	}
}

func TestValidateOTPCode(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
	}{
		{"Valid code", "123456", false},
		{"Too short", "12345", true},
		{"Too long", "1234567", true},
		{"Contains letters", "12a456", true},
		{"Empty string", "", true},
		{"All zeros", "000000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateOTPCode(tt.code)
			hasErrors := errors.HasErrors()

			if tt.expectError {
				assert.True(t, hasErrors, "Expected validation error for %s", tt.code)
			} else {
				assert.False(t, hasErrors, "Expected no validation error for %s", tt.code)
			}
		})
	}
}
