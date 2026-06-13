package service

import (
	"financial-data-aggregator-backend/internal/models"
	"financial-data-aggregator-backend/internal/repository"
)

type UserService interface {
	GetUserProfile(id string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserProfile(id string) (*models.User, error) {
	return s.repo.FindById(id)
}
