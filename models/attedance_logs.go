package models

import "time"

type AttedanceLogs struct {
	ID          int64   `gorm:"primaryKey;autoIncrement"`
	UserID      int64   `gorm:"index"`
	TokenID     *int64  `gorm:"index"` // nullable — manual status (alfa/sakit/izin) tidak perlu token
	Status      string  `gorm:"type:enum('hadir','telat','alfa','sakit','izin');default:'alfa'"`
	CapturedIp  *string `gorm:"type:varchar(45)"`
	ClockInTime time.Time

	User  Users            `gorm:"foreignKey:UserID;references:ID"`
	Token *AttedanceTokens `gorm:"foreignKey:TokenID;references:ID"`
}
