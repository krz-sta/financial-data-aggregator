package main

import (
	"financial-data-aggregator-backend/internal/config"
	"financial-data-aggregator-backend/internal/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var db *gorm.DB
	var err error

	router := gin.Default()

	cfg := config.LoadConfig()
	dsn := cfg.DB.GetDsn()
	rAddr := cfg.Router.GetRouterConfig()

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err.Error())
	}

	fmt.Println("Connected to DB")

	err = db.AutoMigrate(&models.User{}, &models.PortfolioItem{})
	if err != nil {
		log.Fatalf("Couldn't automigrate: %v", err.Error())
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "backen działa",
		})
	})

	router.Run(rAddr)
}
