package services

import (
	"log"
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/repo"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// InitAttendanceCron menginisialisasi cron job untuk pengecekan absensi harian.
// Mengembalikan instance *cron.Cron agar bisa di-Stop() saat graceful shutdown.
func InitAttendanceCron(db *gorm.DB) *cron.Cron {
	c := cron.New()

	// Jalankan setiap menit, cek apakah dalam window waktu notifikasi
	c.AddFunc("* * * * *", func() {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		now := time.Now().In(loc)

		// Ambil settings dari DB
		settings, err := repo.GetNotificationSettingsMap(db)
		if err != nil {
			log.Printf("[CRON] Gagal ambil settings: %v", err)
			return
		}

		// Cek apakah WA enabled
		if settings["wa_enabled"] != "true" {
			return
		}

		// Parse jam mulai dan akhir
		startStr := settings["wa_check_start"]
		endStr := settings["wa_check_end"]

		if startStr == "" || endStr == "" {
			log.Println("[CRON] Jam pengecekan belum diatur.")
			return
		}

		startTime, err := time.ParseInLocation("15:04", startStr, loc)
		if err != nil {
			log.Printf("[CRON] Format wa_check_start salah: %v", err)
			return
		}
		endTime, err := time.ParseInLocation("15:04", endStr, loc)
		if err != nil {
			log.Printf("[CRON] Format wa_check_end salah: %v", err)
			return
		}

		// Set ke tanggal hari ini
		startToday := time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, loc)
		endToday := time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), 0, 0, loc)

		// Cek apakah sekarang dalam range
		if now.Before(startToday) || now.After(endToday) {
			return // di luar jam pengecekan
		}

		log.Println("[CRON] Dalam window pengecekan — memulai notifikasi WA...")
		CheckAndNotifyAbsentStudents(db)
	})

	c.Start()
	log.Println("[CRON] Attendance cron job aktif (interval: setiap 1 menit, cek window dari DB settings).")

	return c
}
