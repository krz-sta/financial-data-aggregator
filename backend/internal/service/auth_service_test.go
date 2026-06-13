package service

import (
	"financial-data-aggregator-backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "test_secret_key")

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("EmailExists", "test@test.com").Return(false, nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()

		user, err := service.Register("test@test.com", "Test User", "password123")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@test.com", user.Email)
		assert.Equal(t, "Test User", user.DisplayName)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		mockRepo.On("EmailExists", "exist@test.com").Return(true, nil).Once()

		user, err := service.Register("exist@test.com", "Exist User", "password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email already exists", err.Error())
	})

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "test_secret_key")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &models.User{
		Email:        "test@test.com",
		PasswordHash: string(hash),
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByEmail", "test@test.com").Return(mockUser, nil).Once()

		token, err := service.Login("test@test.com", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		mockRepo.On("FindByEmail", "test@test.com").Return(mockUser, nil).Once()

		token, err := service.Login("test@test.com", "wrongpassword")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "invalide name or password", err.Error())
	})

	mockRepo.AssertExpectations(t)
}
