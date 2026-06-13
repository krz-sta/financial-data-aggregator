package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPortfolioService_AddItem(t *testing.T) {
	mockRepo := new(MockPortfolioRepository)
	service := NewPortfolioService(mockRepo)

	userID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("AddItem", mock.AnythingOfType("*models.PortfolioItem")).Return(nil).Once()
		err := service.AddItem(userID, "AAPL", 10.5)
		assert.NoError(t, err)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		err := service.AddItem("invalid-uuid", "AAPL", 10.5)
		assert.Error(t, err)
		assert.Equal(t, "wrong id format", err.Error())
	})
	mockRepo.AssertExpectations(t)
}

func TestPortfolioService_DeleteItem(t *testing.T) {
	mockRepo := new(MockPortfolioRepository)
	service := NewPortfolioService(mockRepo)

	userUUID := uuid.New()
	itemUUID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("DeleteItem", userUUID, itemUUID).Return(nil).Once()
		err := service.DeleteItem(userUUID.String(), itemUUID.String())
		assert.NoError(t, err)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		err := service.DeleteItem("invalid", itemUUID.String())
		assert.Error(t, err)
		assert.Equal(t, "incorrect id format", err.Error())
	})

	mockRepo.AssertExpectations(t)
}
