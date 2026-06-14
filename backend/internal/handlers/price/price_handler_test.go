package price

import (
	"context"
	"errors"
	"financial-data-aggregator-backend/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPriceService struct {
	mock.Mock
}

func (m *MockPriceService) StartWorker(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockPriceService) GetRates(ctx context.Context) map[string]float64 {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]float64)
}

func (m *MockPriceService) GetHistory(ctx context.Context, symbol string) ([]models.HistoryPoint, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.HistoryPoint), args.Error(1)
}

func TestGetRates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockPriceService)
	handler := NewHandler(mockService)

	router := gin.New()
	router.GET("/api/rates", handler.GetRates)

	t.Run("Successs", func(t *testing.T) {
		mockRates := map[string]float64{
			"BTC": 50000.0,
			"ETH": 3000.0,
		}

		mockService.On("GetRates", mock.Anything).Return(mockRates).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/rates", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	mockService.AssertExpectations(t)
}

func TestGetHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockPriceService)
	handler := NewHandler(mockService)

	router := gin.New()
	router.GET("/api/rates/history/:symbol", handler.GetHistory)

	t.Run("Success", func(t *testing.T) {
		mockHistory := []models.HistoryPoint{
			{Timestamp: 1620000000, Price: 50000.0},
		}

		mockService.On("GetHistory", mock.Anything, "BTC").Return(mockHistory, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/rates/history/BTC", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService.On("GetHistory", mock.Anything, "INVALID").Return(([]models.HistoryPoint)(nil), errors.New("unsupported asset")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/rates/history/INVALID", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockService.AssertExpectations(t)
}
