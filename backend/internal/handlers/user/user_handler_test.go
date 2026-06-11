package user

import (
	"financial-data-aggregator-backend/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserProfile(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	handler := NewHandler(mockService)

	router := gin.New()
	// Dummy middleware simulating AuthMiddleware behavior for test
	router.POST("/api/protected/profile", func(c *gin.Context) {
		c.Set("userID", "123")
		handler.GetProfile(c)
	})

	t.Run("Success", func(t *testing.T) {
		mockUser := &models.User{Email: "test@test.com", DisplayName: "Test User"}
		mockService.On("GetUserProfile", "123").Return(mockUser, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/protected/profile", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	mockService.AssertExpectations(t)
}
