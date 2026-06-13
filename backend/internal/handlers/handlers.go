package handlers

import (
	"context"
	"financial-data-aggregator-backend/internal/handlers/asset"
	"financial-data-aggregator-backend/internal/handlers/auth"
	"financial-data-aggregator-backend/internal/handlers/health"
	"financial-data-aggregator-backend/internal/handlers/portfolio"
	"financial-data-aggregator-backend/internal/handlers/price"
	"financial-data-aggregator-backend/internal/handlers/user"
	"financial-data-aggregator-backend/internal/middleware"
	"financial-data-aggregator-backend/internal/repository"
	"financial-data-aggregator-backend/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, jwtKey string, redis *redis.Client) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userRepo := repository.NewUserRepository(db)
	portfolioRepo := repository.NewPortfolioRepository(db)

	authService := service.NewAuthService(userRepo, jwtKey)
	userService := service.NewUserService(userRepo)
	portfolioService := service.NewPortfolioService(portfolioRepo)
	assetService := service.NewAssetService()

	priceService := service.NewPriceService(redis, assetService)
	priceService.StartWorker(context.Background())

	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)
	healthHandler := health.NewHandler(db)
	portfolioHandler := portfolio.NewHandler(portfolioService)
	assetHandler := asset.NewHandler(assetService)
	priceHandler := price.NewHandler(priceService)

	api := router.Group("/api")
	{
		api.GET("/assets", assetHandler.GetAssets)
		api.GET("/rates", priceHandler.GetRates)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		protected := api.Group("/protected").Use(middleware.AuthMiddleware(jwtKey))
		{
			protected.POST("/profile", userHandler.GetProfile)

			portfolioGroup := api.Group("/portfolio")
			{
				portfolioGroup.POST("", portfolioHandler.AddPortfolioItem)
				portfolioGroup.DELETE("/:id", portfolioHandler.DeletePortfolioItem)
			}
		}

		healthGroup := api.Group("/health")
		{
			healthGroup.GET("/db", healthHandler.DBHealth)
		}
	}
}
