package service

import (
	"financial-data-aggregator-backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) EmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) FindById(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Mock PortfolioRepository
type MockPortfolioRepository struct {
	mock.Mock
}

func (m *MockPortfolioRepository) AddItem(item *models.PortfolioItem) error {
	return m.Called(item).Error(0)
}

func (m *MockPortfolioRepository) DeleteItem(userUUID, itemUUID uuid.UUID) error {
	return m.Called(userUUID, itemUUID).Error(0)
}
