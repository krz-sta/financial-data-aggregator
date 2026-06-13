package main

import (
	"financial-data-aggregator-backend/internal/config"
	"financial-data-aggregator-backend/internal/database"
	"financial-data-aggregator-backend/internal/handlers"
	"financial-data-aggregator-backend/internal/models"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	router := gin.Default()

	cfg := config.LoadConfig()

	db, err := database.NewPostgres(cfg.DB.GetDsn(), &gorm.Config{})
	if err != nil {
		log.Fatalf("Couldn't connect to databse: %v", err.Error())
	}
	fmt.Println("Connected to db")

	err = db.AutoMigrate(&models.User{}, &models.PortfolioItem{})
	if err != nil {
		log.Fatalf("Couldn't automigrate: %v", err.Error())
	}
	fmt.Println("Automigrated db")

	redis, err := database.NewRedis(cfg.Redis.GetRedisConfig(), cfg.Redis.Password)
	if err != nil {
		log.Fatalf("Couldn't connect to redis: %v", err.Error())
	}
	defer redis.Close()
	fmt.Println("Connected to redis")

	handlers.SetupRoutes(router, db, cfg.JWTKey, redis)

	router.Run(cfg.Router.GetRouterConfig())
}
