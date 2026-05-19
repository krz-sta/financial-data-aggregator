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
			name:           "Password too short",
			payload:        `{"email":"test@test.com", "password":"short"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Wrong email format",
			payload:        `{"email":"zly-email", "password":"supersecret123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "No email",
			payload:        `{"password":"supersecret123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "No password",
			payload:        `{"email":"test@test.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Correct registration",
			payload:        `{"email":"test@test.com", "password":"supersecret123"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Correct registration (user already exists)",
			payload:        `{"email":"test@test.com", "password":"supersecret123"}`,
			expectedStatus: http.StatusConflict,
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
				t.Errorf("%s: expected status %d, recieved %d", tt.name, tt.expectedStatus, w.Code)
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
			name:           "Password too short",
			payload:        `{"email":"test@test.com", "password":"short"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Wrong email format",
			payload:        `{"email":"zly-email", "password":"supersecret123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "No email",
			payload:        `{"password":"supersecret123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "No password",
			payload:        `{"email":"test@test.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unregistered user (mail)",
			payload:        `{"email":"testbad@bad.com", "password":"supersecret123"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong password",
			payload:        `{"email":"test@test.com", "password":"supersecret321bad"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Correct login",
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
				t.Errorf("%s: expected status %d, recieved %d", tt.name, tt.expectedStatus, w.Code)
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

	t.Run("Registration", func(t *testing.T) {
		testRegisterValidation(t, db, jwt_secret)
	})

	t.Run("Login", func(t *testing.T) {
		testLoginValidation(t, db, jwt_secret)
	})
}
