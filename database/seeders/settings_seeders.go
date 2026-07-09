package seeders

import (
	"log"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
)

// SeedNotificationSettings — auto-seed default notification settings jika tabel kosong.
// Ini memastikan fitur WA langsung aktif tanpa harus insert manual.
func SeedNotificationSettings() error {
	var count int64
	database.DB.Model(&models.NotificationSettings{}).Count(&count)

	if count > 0 {
		log.Println("[Seeder] notification_settings sudah terisi, skip seeding.")
		return nil
	}

	defaults := []models.NotificationSettings{
		{
			SettingKey:   "wa_check_start",
			SettingValue: "08:00",
			Description:  "Jam mulai pengecekan notifikasi WA (format HH:MM)",
		},
		{
			SettingKey:   "wa_check_end",
			SettingValue: "09:00",
			Description:  "Jam akhir pengecekan notifikasi WA (format HH:MM)",
		},
		{
			SettingKey:   "wa_enabled",
			SettingValue: "true",
			Description:  "Aktifkan/nonaktifkan notifikasi WA (true/false)",
		},
		{
			SettingKey:   "wa_message_template",
			SettingValue: "Assalamualaikum, kami informasikan bahwa anak Bapak/Ibu *{nama}* (NISN: {nisn}, Kelas: {kelas}) hari ini tercatat *{status}*. Mohon perhatiannya. Terima kasih.",
			Description:  "Template pesan WA",
		},
		{
			SettingKey:   "school_name",
			SettingValue: "SMK PLUS PELITA NUSANTARA",
			Description:  "Nama Sekolah untuk Notifikasi WA",
		},
	}

	for _, s := range defaults {
		if err := database.DB.Create(&s).Error; err != nil {
			log.Printf("[Seeder] Gagal insert setting '%s': %v", s.SettingKey, err)
		} else {
			log.Printf("[Seeder] Inserted setting: %s = %s", s.SettingKey, s.SettingValue)
		}
	}

	log.Println("[Seeder] Default notification_settings berhasil di-seed.")
	return nil
}
