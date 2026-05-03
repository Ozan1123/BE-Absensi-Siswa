package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/repo"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

const (
	StatusAlfa  = "alfa"
	StatusSakit = "sakit"
	StatusIzin  = "izin"
)

// struct ringan buat nampung data target notif
type notifTarget struct {
	UserID      int64
	FullName    string
	Nisn        string
	ClassGroup  string
	ParentPhone string
	Status      string
}

// NormalizePhone — konversi 08xxx jadi 628xxx buat format JID WhatsApp
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

// SendWhatsAppMessage — kirim pesan WA lewat whatsmeow client
func SendWhatsAppMessage(phone, message string) (string, error) {
	if WAClient == nil || !WAClient.IsConnected() {
		return "", fmt.Errorf("WhatsApp client belum terhubung")
	}

	jid := types.NewJID(NormalizePhone(phone), types.DefaultUserServer)

	resp, err := WAClient.SendMessage(context.Background(), jid, &waE2E.Message{
		Conversation: proto.String(message),
	})
	if err != nil {
		return "", fmt.Errorf("gagal kirim pesan: %w", err)
	}

	return fmt.Sprintf("sent_at:%s", resp.Timestamp.String()), nil
}

// BuildNotificationMessage — bikin teks pesan dinamis sesuai status (pake switch-case)
func BuildNotificationMessage(nama, nisn, kelas, status string) string {
	header := fmt.Sprintf(
		"Assalamualaikum Wr. Wb.\n\nYth. Bapak/Ibu Orang Tua/Wali dari:\n"+
			"  Nama  : *%s*\n"+
			"  NISN  : %s\n"+
			"  Kelas : %s\n\n",
		nama, nisn, kelas,
	)

	var body string
	switch strings.ToLower(status) {
	case StatusAlfa:
		body = fmt.Sprintf(
			"Kami informasikan bahwa ananda *%s* hari ini terpantau *BELUM MELAKUKAN ABSENSI (ALFA)*. "+
				"Mohon Bapak/Ibu dapat mengonfirmasi kehadiran putra/putri Anda.",
			nama,
		)
	case StatusSakit:
		body = fmt.Sprintf(
			"Kami informasikan bahwa hari ini ananda *%s* tidak dapat mengikuti kegiatan belajar mengajar karena *SAKIT*. "+
				"Kami pihak sekolah mendoakan agar ananda lekas sembuh dan dapat kembali beraktivitas seperti biasa.",
			nama,
		)
	case StatusIzin:
		body = fmt.Sprintf(
			"Kami informasikan bahwa hari ini ananda *%s* tidak dapat mengikuti kegiatan belajar mengajar dengan keterangan *IZIN*. "+
				"Terima kasih kepada Bapak/Ibu atas informasi yang telah disampaikan.",
			nama,
		)
	default:
		body = fmt.Sprintf(
			"Kami informasikan bahwa ananda *%s* hari ini tercatat dengan status *%s*.",
			nama, strings.ToUpper(status),
		)
	}

	footer := "\n\nTerima kasih atas perhatian Bapak/Ibu.\nWassalamualaikum Wr. Wb.\n\n_Pesan ini dikirim secara otomatis oleh Sistem Absensi Sekolah SMK PLUS PELITA NUSANTARA._"

	return header + body + footer
}

// processNotificationBatch — proses kirim notif per-batch, udah include spam guard + rate limit
func processNotificationBatch(db *gorm.DB, targets []notifTarget, today string) (sent, skipped, failed int) {
	for _, t := range targets {
		// spam guard: kalo status ini udah pernah dikirim hari ini, skip aja
		if repo.IsNotificationSentToday(db, t.UserID, t.Status, today) {
			log.Printf("[WA] Skip %s (status: %s) — udah dikirim hari ini.", t.FullName, t.Status)
			skipped++
			continue
		}

		// ga ada nomor ortu? ya skip juga dong
		if t.ParentPhone == "" {
			log.Printf("[WA] Skip %s — nomor ortu kosong.", t.FullName)
			skipped++
			continue
		}

		// bangun pesan sesuai status
		message := BuildNotificationMessage(t.FullName, t.Nisn, t.ClassGroup, t.Status)

		log.Printf("[WA] Ngirim ke %s buat siswa %s (status: %s)...",
			NormalizePhone(t.ParentPhone), t.FullName, t.Status)

		responseStatus, err := SendWhatsAppMessage(t.ParentPhone, message)

		deliveryStatus := "success"
		if err != nil {
			log.Printf("[WA] Gagal kirim ke %s: %v", t.ParentPhone, err)
			deliveryStatus = "failed"
			failed++
		} else {
			log.Printf("[WA] Sukses kirim ke %s buat %s", t.ParentPhone, t.FullName)
			sent++
		}

		// catat ke tabel notification_logs
		db.Create(&models.NotificationLogs{
			UserID:         t.UserID,
			Phone:          NormalizePhone(t.ParentPhone),
			Status:         t.Status,
			Message:        message,
			SentDate:       today,
			ResponseStatus: deliveryStatus + ": " + responseStatus,
		})

		// rate limit 2 detik biar ga di-ban Meta
		time.Sleep(2 * time.Second)
	}

	return sent, skipped, failed
}

// CheckAndNotifyAllStatuses — fungsi utama yg dipanggil cron, ngurus semua status sekaligus
func CheckAndNotifyAllStatuses(db *gorm.DB) {
	settings, err := repo.GetNotificationSettingsMap(db)
	if err != nil {
		log.Printf("[WA] Gagal ambil settings: %v", err)
		return
	}

	if settings["wa_enabled"] != "true" {
		log.Println("[WA] Notifikasi WA lagi off.")
		return
	}

	if WAClient == nil || !WAClient.IsConnected() {
		log.Println("[WA] Client belum konek, skip dulu.")
		return
	}

	today := repo.TodayDateString()
	var allTargets []notifTarget

	// 1) Tarik data siswa ALFA (belum absen sama sekali)
	alfaStudents, err := repo.GetUnattendedStudents(db)
	if err != nil {
		log.Printf("[WA] Gagal ambil data ALFA: %v", err)
	} else {
		for _, s := range alfaStudents {
			allTargets = append(allTargets, notifTarget{
				UserID: s.ID, FullName: s.FullName, Nisn: s.Nisn,
				ClassGroup: s.ClassGroup, ParentPhone: s.ParentPhone,
				Status: StatusAlfa,
			})
		}
		log.Printf("[WA] ALFA: %d siswa.", len(alfaStudents))
	}

	// 2) Tarik data siswa SAKIT dan IZIN sekaligus (1 query)
	sakitIzinStudents, err := repo.GetStudentsByStatusToday(db, []string{StatusSakit, StatusIzin})
	if err != nil {
		log.Printf("[WA] Gagal ambil data SAKIT/IZIN: %v", err)
	} else {
		for _, s := range sakitIzinStudents {
			allTargets = append(allTargets, notifTarget{
				UserID: s.ID, FullName: s.FullName, Nisn: s.Nisn,
				ClassGroup: s.ClassGroup, ParentPhone: s.ParentPhone,
				Status: s.Status, // Status dinamis dapet dari DB (sakit / izin)
			})
		}
		log.Printf("[WA] SAKIT/IZIN: %d siswa.", len(sakitIzinStudents))
	}

	if len(allTargets) == 0 {
		log.Println("[WA] Ga ada siswa yg perlu dinotif hari ini.")
		return
	}

	log.Printf("[WA] Total %d siswa masuk antrian notif.", len(allTargets))

	sent, skipped, failed := processNotificationBatch(db, allTargets, today)
	log.Printf("[WA] Done — Terkirim: %d | Skip: %d | Gagal: %d", sent, skipped, failed)
}

// CheckAndNotifyAbsentStudents — alias biar cron_service.go ga perlu diubah
func CheckAndNotifyAbsentStudents(db *gorm.DB) {
	CheckAndNotifyAllStatuses(db)
}

// TestSendWhatsApp — kirim pesan test ke nomor tertentu (buat debugging)
func TestSendWhatsApp(phone, message string) (string, error) {
	if WAClient == nil || !WAClient.IsConnected() {
		return "", fmt.Errorf("WhatsApp client belum terhubung. Pastikan sudah scan QR")
	}
	return SendWhatsAppMessage(phone, message)
}
