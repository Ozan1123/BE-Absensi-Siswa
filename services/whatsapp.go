package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/repo"
	"gorm.io/gorm"
)

// NormalizePhone mengubah format 08xxx menjadi 628xxx untuk Fonnte API
func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if strings.HasPrefix(phone, "08") {
		return "62" + phone[1:]
	}
	if strings.HasPrefix(phone, "+62") {
		return phone[1:] // hapus "+" saja
	}
	return phone
}

// SendWhatsAppMessage mengirim pesan WA via Fonnte API
func SendWhatsAppMessage(apiKey, phone, message string) (string, error) {
	apiURL := "https://api.fonnte.com/send"

	normalizedPhone := NormalizePhone(phone)

	data := url.Values{}
	data.Set("target", normalizedPhone)
	data.Set("message", message)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("gagal membuat request: %w", err)
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal mengirim request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	responseStr := string(body)

	if resp.StatusCode != 200 {
		return responseStr, fmt.Errorf("fonnte API error (status %d): %s", resp.StatusCode, responseStr)
	}

	return responseStr, nil
}

// FormatNotificationMessage mengganti placeholder di template pesan
func FormatNotificationMessage(template, nama, nisn, kelas, status string) string {
	msg := template
	msg = strings.ReplaceAll(msg, "{nama}", nama)
	msg = strings.ReplaceAll(msg, "{nisn}", nisn)
	msg = strings.ReplaceAll(msg, "{kelas}", kelas)
	msg = strings.ReplaceAll(msg, "{status}", status)
	return msg
}

// CheckAndNotifyAbsentStudents adalah logic utama pengecekan dan pengiriman notifikasi WA
func CheckAndNotifyAbsentStudents(db *gorm.DB) {
	// Ambil semua settings
	settings, err := repo.GetNotificationSettingsMap(db)
	if err != nil {
		log.Printf("[WA] Gagal mengambil settings: %v", err)
		return
	}

	// Cek apakah notifikasi diaktifkan
	if settings["wa_enabled"] != "true" {
		log.Println("[WA] Notifikasi WA dinonaktifkan.")
		return
	}

	// Ambil API Key dari settings ATAU dari env
	apiKey := settings["wa_api_key"]
	if apiKey == "" {
		log.Println("[WA] API Key Fonnte belum di-set. Skip notifikasi.")
		return
	}

	// Ambil template pesan
	template := settings["wa_message_template"]
	if template == "" {
		template = "Assalamualaikum, kami informasikan bahwa anak Bapak/Ibu *{nama}* (NISN: {nisn}, Kelas: {kelas}) hari ini tercatat *{status}*. Mohon perhatiannya. Terima kasih."
	}

	// Ambil siswa yang belum absen (alfa) hari ini
	absentStudents, err := repo.GetUnattendedStudents(db)
	if err != nil {
		log.Printf("[WA] Gagal mengambil data siswa belum absen: %v", err)
		return
	}

	// Ambil siswa yang sakit/izin hari ini
	sickLeaveStudents, err := repo.GetStudentsByStatusToday(db, []string{"sakit", "izin"})
	if err != nil {
		log.Printf("[WA] Gagal mengambil data siswa sakit/izin: %v", err)
		return
	}

	// Gabungkan semua siswa yang perlu dinotifikasi
	type notifTarget struct {
		UserID      int64
		FullName    string
		Nisn        string
		ClassGroup  string
		ParentPhone string
		Status      string
	}

	var targets []notifTarget

	for _, s := range absentStudents {
		targets = append(targets, notifTarget{
			UserID:      s.ID,
			FullName:    s.FullName,
			Nisn:        s.Nisn,
			ClassGroup:  s.ClassGroup,
			ParentPhone: s.ParentPhone,
			Status:      "alfa",
		})
	}

	for _, s := range sickLeaveStudents {
		targets = append(targets, notifTarget{
			UserID:      s.ID,
			FullName:    s.FullName,
			Nisn:        s.Nisn,
			ClassGroup:  s.ClassGroup,
			ParentPhone: s.ParentPhone,
			Status:      s.Status,
		})
	}

	if len(targets) == 0 {
		log.Println("[WA] Tidak ada siswa yang perlu dinotifikasi.")
		return
	}

	log.Printf("[WA] Ditemukan %d siswa perlu dinotifikasi.", len(targets))

	today := repo.TodayDateString()

	for _, t := range targets {
		// Cek apakah sudah dikirim hari ini
		if repo.IsNotificationSentToday(db, t.UserID, today) {
			log.Printf("[WA] Skip %s — sudah dikirm hari ini.", t.FullName)
			continue
		}

		if t.ParentPhone == "" {
			log.Printf("[WA] Skip %s — tidak ada nomor orang tua.", t.FullName)
			continue
		}

		// Format pesan
		message := FormatNotificationMessage(template, t.FullName, t.Nisn, t.ClassGroup, t.Status)

		// Kirim WA
		log.Printf("[WA] Mengirim ke %s (%s) untuk siswa %s...", t.ParentPhone, NormalizePhone(t.ParentPhone), t.FullName)
		responseStatus, err := SendWhatsAppMessage(apiKey, t.ParentPhone, message)

		status := "success"
		if err != nil {
			log.Printf("[WA] Gagal kirim ke %s: %v", t.ParentPhone, err)
			status = "failed"
		} else {
			log.Printf("[WA] Berhasil kirim ke %s untuk %s", t.ParentPhone, t.FullName)
		}

		// Simpan log
		notifLog := models.NotificationLogs{
			UserID:         t.UserID,
			Phone:          NormalizePhone(t.ParentPhone),
			Status:         t.Status,
			Message:        message,
			SentDate:       today,
			ResponseStatus: status + ": " + responseStatus,
		}
		if err := db.Create(&notifLog).Error; err != nil {
			log.Printf("[WA] Gagal simpan log untuk %s: %v", t.FullName, err)
		}
	}
}

// TestSendWhatsApp mengirim pesan test ke nomor tertentu
func TestSendWhatsApp(phone, message string) (string, error) {
	settings, err := repo.GetNotificationSettingsMap(database.DB)
	if err != nil {
		return "", fmt.Errorf("gagal ambil settings: %w", err)
	}

	apiKey := settings["wa_api_key"]
	if apiKey == "" {
		return "", fmt.Errorf("API Key Fonnte belum diatur. Silakan set di notification settings")
	}

	return SendWhatsAppMessage(apiKey, phone, message)
}
