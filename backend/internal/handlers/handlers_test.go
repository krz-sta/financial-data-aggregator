package handlers

import (
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stdout)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	jwt_secret := "test_secret_key"

	testRedis := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		MaxRetries:  -1,
		DialTimeout: time.Millisecond,
	})

	SetupRoutes(router, nil, jwt_secret, testRedis)

	routes := router.Routes()

	expectedRoutes := map[string]string{
		"/api/auth/register":           "POST",
		"/api/auth/login":              "POST",
		"/api/protected/profile":       "POST",
		"/api/health/db":               "GET",
		"/api/health/redis":            "GET",
		"/api/assets":                  "GET",
		"/api/rates":                   "GET",
		"/api/rates/history/:symbol":   "GET",
		"/api/protected/portfolio":     "POST",
		"/api/protected/portfolio/:id": "DELETE",
	}

	for path, method := range expectedRoutes {
		found := false
		for _, r := range routes {
			if r.Path == path && r.Method == method {
				found = true
				break
			}
		}
		assert.True(t, found, "Missing route: "+method+" "+path)
	}
}
