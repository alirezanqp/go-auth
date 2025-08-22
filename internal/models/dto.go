package models

import "github.com/google/uuid"

type SendOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" validate:"required"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" validate:"required"`
	Code        string `json:"code" binding:"required" validate:"required"`
}

type SendOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type VerifyOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   string    `json:"created_at"`
}

type UsersListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
