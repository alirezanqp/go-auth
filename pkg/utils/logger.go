package utils

import (
	"context"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)

	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	env := os.Getenv("GIN_MODE")
	if env == "release" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}
}

func LogWithContext(ctx context.Context) *logrus.Entry {
	entry := Logger.WithFields(logrus.Fields{})

	if requestID := ctx.Value("request_id"); requestID != nil {
		entry = entry.WithField("request_id", requestID)
	}

	if userID := ctx.Value("user_id"); userID != nil {
		entry = entry.WithField("user_id", userID)
	}

	return entry
}

func LogWithFields(fields map[string]interface{}) *logrus.Entry {
	return Logger.WithFields(fields)
}

func LogRequest(method, path, userAgent, clientIP string, statusCode int, latency time.Duration) {
	Logger.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"user_agent":  userAgent,
		"client_ip":   clientIP,
		"status_code": statusCode,
		"latency_ms":  latency.Milliseconds(),
		"type":        "http_request",
	}).Info("HTTP Request")
}

func LogError(err error, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	fields["error"] = err.Error()
	fields["type"] = "error"

	if pc, file, line, ok := runtime.Caller(1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			fields["function"] = fn.Name()
		}
		fields["file"] = file
		fields["line"] = line
	}

	Logger.WithFields(fields).Error(message)
}

func LogOTPGenerated(phoneNumber, code string, expiresAt time.Time) {
	Logger.WithFields(logrus.Fields{
		"phone_number": maskPhoneNumber(phoneNumber),
		"otp_code":     code,
		"expires_at":   expiresAt.Format(time.RFC3339),
		"type":         "otp_generated",
		"action":       "send_otp",
	}).Info("OTP Generated")
}

func LogOTPVerification(phoneNumber, code string, success bool, reason string) {
	Logger.WithFields(logrus.Fields{
		"phone_number": maskPhoneNumber(phoneNumber),
		"success":      success,
		"reason":       reason,
		"type":         "otp_verification",
		"action":       "verify_otp",
	}).Info("OTP Verification Attempt")
}

func LogUserRegistration(userID, phoneNumber string) {
	Logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"phone_number": maskPhoneNumber(phoneNumber),
		"type":         "user_registration",
		"action":       "register",
	}).Info("New User Registered")
}

func LogUserLogin(userID, phoneNumber string) {
	Logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"phone_number": maskPhoneNumber(phoneNumber),
		"type":         "user_login",
		"action":       "login",
	}).Info("User Login")
}

func LogSecurityEvent(eventType, userID, phoneNumber, details string) {
	Logger.WithFields(logrus.Fields{
		"event_type":   eventType,
		"user_id":      userID,
		"phone_number": maskPhoneNumber(phoneNumber),
		"details":      details,
		"type":         "security",
		"severity":     "warning",
	}).Warn("Security Event")
}

func LogRateLimit(phoneNumber string, attempts, maxAttempts int) {
	Logger.WithFields(logrus.Fields{
		"phone_number": maskPhoneNumber(phoneNumber),
		"attempts":     attempts,
		"max_attempts": maxAttempts,
		"type":         "rate_limit",
		"action":       "blocked",
	}).Warn("Rate Limit Exceeded")
}

func LogDatabaseOperation(operation, table string, success bool, errorMsg string) {
	fields := logrus.Fields{
		"operation": operation,
		"table":     table,
		"success":   success,
		"type":      "database",
	}

	if !success && errorMsg != "" {
		fields["error"] = errorMsg
		Logger.WithFields(fields).Error("Database Operation Failed")
	} else {
		Logger.WithFields(fields).Debug("Database Operation")
	}
}

func maskPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) <= 4 {
		return strings.Repeat("*", len(phoneNumber))
	}

	return phoneNumber[:2] + strings.Repeat("*", len(phoneNumber)-4) + phoneNumber[len(phoneNumber)-2:]
}
