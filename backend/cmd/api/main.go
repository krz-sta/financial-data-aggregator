package main

import (
	"financial-data-aggregator-backend/internal/config"
	"financial-data-aggregator-backend/internal/database"
	"financial-data-aggregator-backend/internal/handlers"
	"financial-data-aggregator-backend/internal/models"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	router := gin.Default()

	cfg := config.LoadConfig()
	dsn := cfg.DB.GetDsn()
	rAddr := cfg.Router.GetRouterConfig()

	db, err := database.NewPostgres(dsn, &gorm.Config{})

	err = db.AutoMigrate(&models.User{}, &models.PortfolioItem{})
	if err != nil {
		log.Fatalf("Couldn't automigrate: %v", err.Error())
	}

	handlers.SetupRoutes(router, db, cfg.JWTKey)

	router.Run(rAddr)
}
