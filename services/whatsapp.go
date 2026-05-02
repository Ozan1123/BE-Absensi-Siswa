package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/repo"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// NormalizePhone mengubah format 08xxx menjadi 628xxx untuk WhatsApp JID
func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if strings.HasPrefix(phone, "08") {
		return "62" + phone[1:]
	}
	if strings.HasPrefix(phone, "+62") {
		return phone[1:]
	}
	return phone
}

// SendWhatsAppMessage mengirim pesan WA via whatsmeow client
func SendWhatsAppMessage(phone, message string) (string, error) {
	if WAClient == nil || !WAClient.IsConnected() {
		return "", fmt.Errorf("WhatsApp client belum terhubung")
	}

	normalizedPhone := NormalizePhone(phone)
	jid := types.NewJID(normalizedPhone, types.DefaultUserServer)

	msg := &waE2E.Message{
		Conversation: proto.String(message),
	}

	resp, err := WAClient.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return "", fmt.Errorf("gagal kirim pesan: %w", err)
	}

	return fmt.Sprintf("sent_at:%s", resp.Timestamp.String()), nil
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

// CheckAndNotifyAbsentStudents — logic utama pengecekan dan pengiriman notifikasi WA
func CheckAndNotifyAbsentStudents(db *gorm.DB) {
	settings, err := repo.GetNotificationSettingsMap(db)
	if err != nil {
		log.Printf("[WA] Gagal mengambil settings: %v", err)
		return
	}

	if settings["wa_enabled"] != "true" {
		log.Println("[WA] Notifikasi WA dinonaktifkan.")
		return
	}

	if WAClient == nil || !WAClient.IsConnected() {
		log.Println("[WA] Client belum terhubung. Skip notifikasi.")
		return
	}

	template := settings["wa_message_template"]
	if template == "" {
		template = "Assalamualaikum, kami informasikan bahwa anak Bapak/Ibu *{nama}* (NISN: {nisn}, Kelas: {kelas}) hari ini tercatat *{status}*. Mohon perhatiannya. Terima kasih."
	}

	absentStudents, err := repo.GetUnattendedStudents(db)
	if err != nil {
		log.Printf("[WA] Gagal mengambil data siswa belum absen: %v", err)
		return
	}

	sickLeaveStudents, err := repo.GetStudentsByStatusToday(db, []string{"sakit", "izin"})
	if err != nil {
		log.Printf("[WA] Gagal mengambil data siswa sakit/izin: %v", err)
		return
	}

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
			UserID: s.ID, FullName: s.FullName, Nisn: s.Nisn,
			ClassGroup: s.ClassGroup, ParentPhone: s.ParentPhone, Status: "alfa",
		})
	}

	for _, s := range sickLeaveStudents {
		targets = append(targets, notifTarget{
			UserID: s.ID, FullName: s.FullName, Nisn: s.Nisn,
			ClassGroup: s.ClassGroup, ParentPhone: s.ParentPhone, Status: s.Status,
		})
	}

	if len(targets) == 0 {
		log.Println("[WA] Tidak ada siswa yang perlu dinotifikasi.")
		return
	}

	log.Printf("[WA] Ditemukan %d siswa perlu dinotifikasi.", len(targets))

	today := repo.TodayDateString()

	for _, t := range targets {
		if repo.IsNotificationSentToday(db, t.UserID, today) {
			log.Printf("[WA] Skip %s — sudah dikirim hari ini.", t.FullName)
			continue
		}

		if t.ParentPhone == "" {
			log.Printf("[WA] Skip %s — tidak ada nomor orang tua.", t.FullName)
			continue
		}

		message := FormatNotificationMessage(template, t.FullName, t.Nisn, t.ClassGroup, t.Status)

		log.Printf("[WA] Mengirim ke %s untuk siswa %s...", NormalizePhone(t.ParentPhone), t.FullName)
		responseStatus, err := SendWhatsAppMessage(t.ParentPhone, message)

		status := "success"
		if err != nil {
			log.Printf("[WA] Gagal kirim ke %s: %v", t.ParentPhone, err)
			status = "failed"
		} else {
			log.Printf("[WA] Berhasil kirim ke %s untuk %s", t.ParentPhone, t.FullName)
		}

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
	if WAClient == nil || !WAClient.IsConnected() {
		return "", fmt.Errorf("WhatsApp client belum terhubung. Pastikan sudah scan QR")
	}
	return SendWhatsAppMessage(phone, message)
}
