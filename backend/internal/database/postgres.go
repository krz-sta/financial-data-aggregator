package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(dsn string, config *gorm.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	db, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err.Error())
	}

	fmt.Println("Connected to DB")

	return db, err
}
