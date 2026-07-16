package models

import "time"

type AdminNotifications struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Judul     string     `gorm:"type:varchar(255)" json:"judul"`
	Pesan     string     `gorm:"type:text" json:"pesan"`
	Tipe      string     `gorm:"type:varchar(50)" json:"tipe"` // "WA error" or "Rekap"
	IsRead    bool       `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
