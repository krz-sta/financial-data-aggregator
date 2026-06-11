package service

import (
	"financial-data-aggregator-backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	mockUser := &models.User{Email: "test@test.com", DisplayName: "Test User"}

	// Oczekujemy wywołania FindById z ID "123" i zwracamy mockUser
	mockRepo.On("FindById", "123").Return(mockUser, nil)

	user, err := userService.GetUserProfile("123")

	assert.NoError(t, err)
	assert.Equal(t, "test@test.com", user.Email)
	mockRepo.AssertExpectations(t)
}
