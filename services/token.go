package services

import (
	"log"
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
)

// StartTokenCleaner — background loop yang berjalan tiap 1 menit.
// Melakukan dua tugas:
//  1. Nonaktifkan token yang sudah melewati valid_until (is_active = false).
//  2. Deteksi token expired yang belum diproses notifikasinya (notification_processed = false),
//     lalu trigger pengiriman notifikasi WA sesuai kategori token dan tandai sebagai processed.
//
// Pendekatan ini menggantikan time.AfterFunc di memori, sehingga tahan terhadap server restart.
func StartTokenCleaner() {
	go func() {
		for {
			now := time.Now()

			// 1. Nonaktifkan token yang sudah expired
			database.DB.
				Model(&models.AttedanceTokens{}).
				Where("is_active = ? AND valid_until < ?", true, now).
				Update("is_active", false)

			// 2. Cari token expired yang belum diproses notifikasinya
			var unprocessedTokens []models.AttedanceTokens
			database.DB.
				Where("is_active = ? AND valid_until < ? AND notification_processed = ?", false, now, false).
				Find(&unprocessedTokens)

			for _, token := range unprocessedTokens {
				log.Printf("[WA-CLEANER] Token %s (kategori: %s) expired & belum diproses — mendelegasikan broadcast notifikasi...",
					token.TokenCode, token.Category)

				// Jalankan di background goroutine agar tidak menyumbat loop utama cleaner
				go func(t models.AttedanceTokens) {
					if t.Category == "hadir" {
						NotifyPresentStudents(database.DB)
					} else if t.Category == "telat" {
						AutoAlfaAndNotify(database.DB)
					}
				}(token)

				// Tandai token sebagai sudah diproses notifikasinya
				database.DB.
					Model(&models.AttedanceTokens{}).
					Where("id = ?", token.ID).
					Update("notification_processed", true)

				log.Printf("[WA-CLEANER] Token %s (kategori: %s) — broadcast notifikasi didelegasikan, ditandai processed.",
					token.TokenCode, token.Category)
			}

			time.Sleep(1 * time.Minute) // cek tiap 1 menit
		}
	}()
}