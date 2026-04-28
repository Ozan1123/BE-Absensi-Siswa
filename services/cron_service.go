package services

import (
	"log"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/repo"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// InitAttendanceCron menginisialisasi cron job untuk pengecekan absensi harian.
// Mengembalikan instance *cron.Cron agar bisa di-Stop() saat graceful shutdown.
func InitAttendanceCron(db *gorm.DB) *cron.Cron {
	c := cron.New()

	// Testing: jalankan setiap 1 menit
	// Produksi: ganti ke "30 8 * * *" (setiap hari jam 08:30)
	c.AddFunc("* * * * *", func() {
		log.Println("[CRON] Memulai pengecekan siswa yang belum absen...")

		students, err := repo.GetUnattendedStudents(db)
		if err != nil {
			log.Printf("[CRON] Gagal mengambil data siswa: %v", err)
			return
		}

		if len(students) == 0 {
			log.Println("[CRON] Semua siswa sudah absen hari ini.")
			return
		}

		log.Printf("[CRON] Ditemukan %d siswa belum absen.", len(students))

		for _, s := range students {
			log.Printf("[CRON-TEST] Akan mengirim WA ke %s untuk siswa %s (NISN: %s) karena belum absen hari ini.",
				s.ParentPhone, s.FullName, s.Nisn)
		}
	})

	c.Start()
	log.Println("[CRON] Attendance cron job aktif (interval: setiap 1 menit).")

	return c
}
