package models

import "time"

type AttedanceTokens struct {
	ID         int64
	TokenCode  string
	CreatedBy  int64
	IsActive   bool
	ValidUntil time.Time
	CreatedAt  time.Time

	User Users `gorm:"foreignKey:CreatedBy;references:ID"`
}
