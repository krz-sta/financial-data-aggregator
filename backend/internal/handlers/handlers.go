package handlers

import (
	"financial-data-aggregator-backend/internal/handlers/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, jwtKey string) {
	authHandler := auth.NewHandler(db, jwtKey)

	api := router.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}
	}
}
