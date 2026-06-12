package models

import "github.com/google/uuid"

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"displayName" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
}
