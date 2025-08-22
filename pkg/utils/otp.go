package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {
	min := int64(100000)
	max := int64(999999)

	n, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	otp := n.Int64() + min
	return fmt.Sprintf("%06d", otp), nil
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}

	return string(result)
}
