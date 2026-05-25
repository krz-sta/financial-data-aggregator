package handlers

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	jwt_secret := "test_secret_key"

	SetupRoutes(router, nil, jwt_secret)

	routes := router.Routes()

	foundRegister := false

	for _, r := range routes {
		if r.Path == "/api/auth/register" && r.Method == "POST" {
			foundRegister = true
		}

		if r.Path == "/api/auth/login" && r.Method == "POST" {
			foundRegister = true
		}

		if r.Path == "/api/protected/profile" && r.Method == "POST" {
			foundRegister = true
		}

		if r.Path == "/api/health/db" && r.Method == "GET" {
			foundRegister = true
		}
	}

	if !foundRegister {
		t.Errorf("Router did not register route")
	}
}
