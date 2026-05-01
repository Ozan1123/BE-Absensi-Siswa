package models

import "time"

type AttedanceTokens struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	TokenCode  string    `gorm:"type:varchar(10);uniqueIndex;not null"`
	CreatedBy  int64     `gorm:"not null"`
	IsActive   bool      `gorm:"type:boolean"`
	LateAfter  time.Time `gorm:"type:datetime"`
	ValidUntil time.Time `gorm:"type:datetime"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	User Users `gorm:"foreignKey:CreatedBy;references:ID"`
}
