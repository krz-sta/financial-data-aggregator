package service

import (
	"errors"
	"financial-data-aggregator-backend/internal/models"
	"financial-data-aggregator-backend/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, name, password string) (*models.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	repo   repository.UserRepository
	jwtKey string
}

func NewAuthService(repo repository.UserRepository, jwtKey string) AuthService {
	return &authService{repo: repo, jwtKey: jwtKey}
}

func (s *authService) Register(email, name, password string) (*models.User, error) {
	exists, err := s.repo.EmailExists(email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		DisplayName:  name,
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalide name or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalide name or password")
	}

	// sub (Subject) - ID użytkownika
	// exp (Expiration Time) - czas wygaśnięcia tokenu w formacie Unix
	// iat (Issued at) - czas wydania tokenu w formacie Unix
	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtKey))
}
