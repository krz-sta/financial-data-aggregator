package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"financial-data-aggregator-backend/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(email, name, password string) (*models.User, error) {
	args := m.Called(email, name, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	handler := NewHandler(mockService)

	router := gin.New()
	router.POST("/api/auth/register", handler.Register)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.RegisterInput{Email: "test@test.com", Name: "Test User", Password: "password123"}
		jsonBody, _ := json.Marshal(reqBody)

		mockUser := &models.User{Email: "test@test.com", DisplayName: "Test User"}
		mockService.On("Register", "test@test.com", "Test User", "password123").Return(mockUser, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(`{"invalid"}`)))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		reqBody := models.RegisterInput{Email: "error@test.com", Name: "Error User", Password: "password123"}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("Register", "error@test.com", "Error User", "password123").Return((*models.User)(nil), errors.New("db error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockService.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	handler := NewHandler(mockService)

	router := gin.New()
	router.POST("/api/auth/login", handler.Login)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.LoginInput{Email: "test@test.com", Password: "password123"}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("Login", "test@test.com", "password123").Return("test_token", nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer([]byte(`{"invalid"}`)))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		reqBody := models.LoginInput{Email: "wrong@test.com", Password: "wrongpassword"}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("Login", "wrong@test.com", "wrongpassword").Return("", errors.New("invalid auth")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	mockService.AssertExpectations(t)
}
