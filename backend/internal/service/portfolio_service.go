package service

import (
	"errors"
	"financial-data-aggregator-backend/internal/models"
	"financial-data-aggregator-backend/internal/repository"

	"github.com/google/uuid"
)

type PortfolioService interface {
	AddItem(userID, symbol string, amount float64) error
	DeleteItem(userID, itemID string) error
}

type portfolioService struct {
	repo repository.PortfolioRepository
}

func NewPortfolioService(repo repository.PortfolioRepository) PortfolioService {
	return &portfolioService{repo: repo}
}

func (s *portfolioService) AddItem(userID, symbol string, amount float64) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("wrong id format")
	}

	newItem := &models.PortfolioItem{
		ID:     uuid.New(),
		UserID: userUUID,
		Symbol: symbol,
		Amount: amount,
	}

	return s.repo.AddItem(newItem)
}

func (s *portfolioService) DeleteItem(userID, itemID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("incorrect id format")
	}

	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		return errors.New("incorrect id format")
	}

	return s.repo.DeleteItem(userUUID, itemUUID)
}
