package models

import "time"

// Note: Untuk menghindari race condition, tabel ini membutuhkan unique index di database:
// ALTER TABLE attedance_logs ADD UNIQUE INDEX uq_user_daily (user_id, (DATE(clock_in_time)));
type AttedanceLogs struct {
	ID          int64   `gorm:"primaryKey;autoIncrement"`
	UserID      int64   `gorm:"index"`
	TokenID     *int64  `gorm:"index"` // nullable — manual status (alfa/sakit) tidak perlu token
	Status      string  `gorm:"type:enum('hadir','telat','alfa','sakit');default:'alfa'"`
	CapturedIp  *string `gorm:"type:varchar(45)"`
	ClockInTime time.Time

	User  Users            `gorm:"foreignKey:UserID;references:ID"`
	Token *AttedanceTokens `gorm:"foreignKey:TokenID;references:ID"`
}
