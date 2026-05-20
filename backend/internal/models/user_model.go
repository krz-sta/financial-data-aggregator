package models

import (
	"time"

	"github.com/google/uuid"
)

type PortfolioItem struct {
	ID        uuid.UUID `gorm:"type:uuid;noy null;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userID"`
	Symbol    string    `gorm:"type:varchar(20);not null;index" json:"symbol"`
	Amount    float64   `gorm:"type:numeric;not null;" json:"amount"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"craetedAt"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"UpdatedAt"`
}

type User struct {
	ID           uuid.UUID       `gorm:"type:uuid;not null;primaryKey" json:"id"`
	Email        string          `gorm:"type:varchar(255);not null" json:"email"`
	PasswordHash string          `gorm:"type:varchar(255);not null" json:"-"`
	DisplayName  string          `gorm:"type:varchar(255); not null" json:"displayName"`
	CreatedAt    time.Time       `gorm:"type:timestamp" json:"createdAt"`
	Portfolio    []PortfolioItem `gorm:"foreignKey:UserID" json:"portfolio,omitempty"`
}
