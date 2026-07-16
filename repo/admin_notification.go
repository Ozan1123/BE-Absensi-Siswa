package repo

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
)

func InsertNotification(judul, pesan, tipe string) error {
	notif := models.AdminNotifications{
		Judul: judul,
		Pesan: pesan,
		Tipe:  tipe,
	}
	return database.DB.Create(&notif).Error
}

func GetUnreadNotifications() ([]models.AdminNotifications, error) {
	var notifs []models.AdminNotifications
	err := database.DB.Where("is_read = ?", false).Order("created_at desc").Find(&notifs).Error
	return notifs, err
}

func MarkAsRead(id int64) error {
	return database.DB.Model(&models.AdminNotifications{}).Where("id = ?", id).Update("is_read", true).Error
}
