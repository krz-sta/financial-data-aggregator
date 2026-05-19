package auth

import (
	"bytes"
	"financial-data-aggregator-backend/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testRegisterValidation(t *testing.T, db *gorm.DB, jwtKey string) {
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
	}{
		{
			name:           "Za krótkie hasło",
			payload:        `{"email":"test@test.com", "password":"short"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Zły format email",
			payload:        `{"email":"zly-email", "password":"supersecret123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Brak adresu email",
			payload:        `{"password":"supersecret123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Brak hasła",
			payload:        `{"email":"test@test.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Poprawna rejestracja",
			payload:        `{"email":"test@test.com", "password":"supersecret123"}`,
			expectedStatus: http.StatusCreated,
		},
	}

	gin.SetMode(gin.TestMode)
	h := NewHandler(db, jwtKey)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(tt.payload)))

			h.Register(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: oczekiwano statusu %d, otrzymano %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}

func testLoginValidation(t *testing.T, db *gorm.DB, jwtKey string) {
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
	}{
		{
			name:           "Za krótkie hasło",
			payload:        `{"email":"test@test.com", "password":"short"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Zły format email",
			payload:        `{"email":"zly-email", "password":"supersecret123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Brak adresu email",
			payload:        `{"password":"supersecret123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Brak hasła",
			payload:        `{"email":"test@test.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Poprawna logowanie",
			payload:        `{"email":"test@test.com", "password":"supersecret123"}`,
			expectedStatus: http.StatusOK,
		},
	}

	gin.SetMode(gin.TestMode)
	h := NewHandler(db, jwtKey)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer([]byte(tt.payload)))

			h.Login(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: oczekiwano statusu %d, otrzymano %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuth(t *testing.T) {
	var err error
	var db *gorm.DB

	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	jwt_secret := "test_secret_key"

	t.Run("Rejestracja", func(t *testing.T) {
		testRegisterValidation(t, db, jwt_secret)
	})

	t.Run("Logowanie", func(t *testing.T) {
		testLoginValidation(t, db, jwt_secret)
	})
}
