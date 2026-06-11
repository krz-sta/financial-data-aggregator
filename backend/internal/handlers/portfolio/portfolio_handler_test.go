package portfolio

import (
	"bytes"
	"encoding/json"
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

	mockService.AssertExpectations(t)
}
