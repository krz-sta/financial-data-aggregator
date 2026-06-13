package repository

import (
	"errors"
	"financial-data-aggregator-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PortfolioRepository interface {
	AddItem(item *models.PortfolioItem) error
	DeleteItem(userUUID, itemUUID uuid.UUID) error
}

type portfolioRepository struct {
	db *gorm.DB
}

func NewPortfolioRepository(db *gorm.DB) PortfolioRepository {
	return &portfolioRepository{db: db}
}

func (r *portfolioRepository) AddItem(item *models.PortfolioItem) error {
	return r.db.Create(item).Error
}

func (r *portfolioRepository) DeleteItem(userUUID, itemUUID uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", itemUUID, userUUID).Delete(&models.PortfolioItem{})
	if result.Error != nil {
		return errors.New("couldn't delete item")
	}

	if result.RowsAffected == 0 {
		return errors.New("couldn't find the item")
	}

	return nil
}
