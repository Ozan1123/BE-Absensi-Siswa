package services

import (
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
)

func StartTokenCleaner() {
	go func() {
		for {
			database.DB.
				Model(&models.AttedanceTokens{}).
				Where("is_active = ? AND valid_until < ?", true, time.Now()).
				Update("is_active", false)

			time.Sleep(1 * time.Minute) // cek tiap 1 menit
		}
	}()
}