package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(dsn string, config *gorm.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	db, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}
	return db, err
}
