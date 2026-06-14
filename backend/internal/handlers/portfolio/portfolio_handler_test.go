package portfolio

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPortfolioService struct {
	mock.Mock
}

func (m *MockPortfolioService) AddItem(userID, symbol string, amount float64) error {
	return m.Called(userID, symbol, amount).Error(0)
}

func (m *MockPortfolioService) DeleteItem(userID, itemID string) error {
	return m.Called(userID, itemID).Error(0)
}

func TestAddPortfolioItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockPortfolioService)
	handler := NewHandler(mockService)

	router := gin.New()
	router.POST("/api/protected/portfolio", func(c *gin.Context) {
		c.Set("userID", "123")
		handler.AddPortfolioItem(c)
	})

	t.Run("Success", func(t *testing.T) {
		reqBody := addInput{Symbol: "AAPL", Amount: 15.5}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("AddItem", "123", "AAPL", 15.5).Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/protected/portfolio", bytes.NewBuffer(jsonBody))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized - No User ID", func(t *testing.T) {
		routerNoUser := gin.New()
		routerNoUser.POST("/api/protected/portfolio", handler.AddPortfolioItem)

		reqBody := addInput{Symbol: "AAPL", Amount: 15.5}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/protected/portfolio", bytes.NewBuffer(jsonBody))
		routerNoUser.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/protected/portfolio", bytes.NewBuffer([]byte(`{"invalid"}`)))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		reqBody := addInput{Symbol: "AAPL", Amount: 15.5}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("AddItem", "123", "AAPL", 15.5).Return(errors.New("db error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/protected/portfolio", bytes.NewBuffer(jsonBody))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockService.AssertExpectations(t)
}

func TestDeletePortfolioItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockPortfolioService)
	handler := NewHandler(mockService)

	router := gin.New()
	router.DELETE("/api/protected/portfolio/:id", func(c *gin.Context) {
		c.Set("userID", "123")
		handler.DeletePortfolioItem(c)
	})

	t.Run("Success", func(t *testing.T) {
		mockService.On("DeleteItem", "123", "456").Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/protected/portfolio/456", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized - No User ID", func(t *testing.T) {
		routerNoUser := gin.New()
		routerNoUser.DELETE("/api/protected/portfolio/:id", handler.DeletePortfolioItem)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/protected/portfolio/456", nil)
		routerNoUser.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService.On("DeleteItem", "123", "456").Return(errors.New("delete error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/protected/portfolio/456", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockService.AssertExpectations(t)
}
