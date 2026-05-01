package models

import "time"

type NotificationSettings struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SettingKey   string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"setting_key"`
	SettingValue string    `gorm:"type:text;not null" json:"setting_value"`
	Description  string    `gorm:"type:varchar(255)" json:"description"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
