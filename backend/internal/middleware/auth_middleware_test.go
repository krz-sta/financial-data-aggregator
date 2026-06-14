package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(userID string, expTime time.Time, secretKey string) string {
	testClaims := jwt.MapClaims{
		"sub": userID,
		"exp": float64(expTime.Unix()),
		"iat": time.Now().Unix(),
	}

	testToken := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	testTokenString, err := testToken.SignedString([]byte(secretKey))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return testTokenString
}

func TestMiddleware(t *testing.T) {
	jwt_secret := "new_test_key"

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
	}{
		{
			name:           "No header",
			payload:        "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong header format",
			payload:        "Bearer" + generateTestToken("testUser", time.Now().Add(time.Hour), jwt_secret),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong token",
			payload:        "Bearer wrongToken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong secret key",
			payload:        "Bearer" + " " + generateTestToken("testUser", time.Now().Add(time.Hour), "wrongPassword"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Expired token",
			payload:        "Bearer" + " " + generateTestToken("testUser", time.Now().Add(-time.Hour), jwt_secret),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Successful auth",
			payload:        "Bearer" + " " + generateTestToken("testUser", time.Now().Add(time.Hour), jwt_secret),
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing sub claim",
			payload: "Bearer " + func() string {
				testClaims := jwt.MapClaims{
					"exp": float64(time.Now().Add(time.Hour).Unix()),
					"iat": time.Now().Unix(),
				}
				testToken := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
				str, _ := testToken.SignedString([]byte(jwt_secret))
				return str
			}(),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(AuthMiddleware(jwt_secret))

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)

			req.Header.Set("Authorization", tt.payload)

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, recieved %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}
