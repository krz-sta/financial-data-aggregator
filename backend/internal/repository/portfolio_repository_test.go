package repository

import (
	"financial-data-aggregator-backend/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPortfolioRepository(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	repo := NewPortfolioRepository(db)

	userID := uuid.New()
	itemID := uuid.New()

	err := db.Create(&models.User{
		ID:           userID,
		Email:        "portfoliotest@test.com",
		DisplayName:  "portfoliotest",
		PasswordHash: "hashedpassword",
	}).Error
	assert.NoError(t, err, "failed to create user for portfolio test")

	t.Run("Add Item", func(t *testing.T) {
		item := &models.PortfolioItem{
			ID:     itemID,
			UserID: userID,
			Symbol: "AAPL",
			Amount: 15.5,
		}

		err := repo.AddItem(item)
		assert.NoError(t, err)
	})

	t.Run("Delete Item Success", func(t *testing.T) {
		err := repo.DeleteItem(userID, itemID)
		assert.NoError(t, err)
	})

	t.Run("Delete Non-existent Item", func(t *testing.T) {
		err := repo.DeleteItem(uuid.New(), uuid.New())
		assert.Error(t, err)
		assert.Equal(t, "couldn't find the item", err.Error())
	})
}
