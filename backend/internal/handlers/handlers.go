package handlers

import (
	"financial-data-aggregator-backend/internal/handlers/auth"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, jwtKey string) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Conten-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Lenght"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
