package services

import (
	"context"
	"fmt"
	"log"
	"regexp"
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

// NormalizePhone — konversi nomor HP ke format JID WhatsApp (628xxx)
// Menghapus semua karakter non-angka (spasi, strip, kurung, dll)
// lalu mengkonversi awalan 08/+62 ke 62.
func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	// Hapus semua karakter non-digit
	re := regexp.MustCompile(`[^0-9]`)
	phone = re.ReplaceAllString(phone, "")
	if strings.HasPrefix(phone, "08") {
		return "62" + phone[1:]
	}
	// Jika dimulai dengan 62 (dari +62 yang sudah di-strip), langsung return
	if strings.HasPrefix(phone, "62") {
		return phone
	}
	return phone
}

// SendWhatsAppMessage — kirim pesan WA lewat whatsmeow client
// Menggunakan context timeout 15 detik agar tidak hang jika koneksi mati.
func SendWhatsAppMessage(phone, message string) (string, error) {
	if WAClient == nil || !WAClient.IsConnected() {
		return "", fmt.Errorf("WhatsApp client belum terhubung")
	}

	jid := types.NewJID(NormalizePhone(phone), types.DefaultUserServer)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := WAClient.SendMessage(ctx, jid, &waE2E.Message{
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
	case "hadir":
		body = fmt.Sprintf(
			"Kami informasikan bahwa ananda *%s* hari ini telah hadir di sekolah dan melakukan absensi dengan status *HADIR*. "+
				"Terima kasih atas perhatian Bapak/Ibu.",
			nama,
		)
	case "telat":
		body = fmt.Sprintf(
			"Kami informasikan bahwa ananda *%s* hari ini hadir di sekolah namun tercatat *TERLAMBAT*. "+
				"Mohon Bapak/Ibu dapat mengingatkan putra/putri Anda untuk datang tepat waktu.",
			nama,
		)
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

// NotifyPresentStudents — kirim notif hanya untuk siswa yang sudah HADIR (dipanggil setelah QR 1 expired)
func NotifyPresentStudents(db *gorm.DB) {
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

	// Tarik data siswa HADIR saja
	hadirStudents, err := repo.GetStudentsByStatusToday(db, []string{"hadir"})
	if err != nil {
		log.Printf("[WA] Gagal ambil data HADIR: %v", err)
	} else {
		for _, s := range hadirStudents {
			allTargets = append(allTargets, notifTarget{
				UserID: s.ID, FullName: s.FullName, Nisn: s.Nisn,
				ClassGroup: s.ClassGroup, ParentPhone: s.ParentPhone,
				Status: "hadir",
			})
		}
		log.Printf("[WA] HADIR: %d siswa.", len(hadirStudents))
	}

	if len(allTargets) == 0 {
		log.Println("[WA] Ga ada siswa yg perlu dinotif hadir.")
		return
	}

	log.Printf("[WA] Total %d siswa masuk antrian notif hadir.", len(allTargets))
	sent, skipped, failed := processNotificationBatch(db, allTargets, today)
	log.Printf("[WA] Done (Hadir) — Terkirim: %d | Skip: %d | Gagal: %d", sent, skipped, failed)
}

// AutoAlfaAndNotify — set alfa untuk siswa tanpa log, lalu kirim notif telat/sakit/izin/alfa (dipanggil setelah QR 2 expired)
func AutoAlfaAndNotify(db *gorm.DB) {
	// 1. Jalankan Auto-Alfa dulu
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	log.Println("[WA] Memulai proses Auto-Alfa...")
	
	// Cari siswa yang belum absen hari ini
	unattended, err := repo.GetUnattendedStudents(db)
	if err != nil {
		log.Printf("[WA] Gagal ambil siswa tanpa absen untuk auto-alfa: %v", err)
	} else {
		for _, s := range unattended {
			logAbsensi := models.AttedanceLogs{
				UserID:      s.ID,
				Status:      StatusAlfa,
				ClockInTime: now,
			}
			if err := db.Create(&logAbsensi).Error; err != nil {
				log.Printf("[WA] Gagal set alfa untuk siswa %d: %v", s.ID, err)
			}
		}
		log.Printf("[WA] Auto-Alfa selesai: %d siswa ditandai ALFA.", len(unattended))
	}

	// 2. Lanjut ke proses pengiriman notifikasi WA
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

	// Tarik data siswa ALFA, SAKIT, IZIN, TELAT (semua kecuali hadir)
	targetStudents, err := repo.GetStudentsByStatusToday(db, []string{StatusAlfa, StatusSakit, StatusIzin, "telat"})
	if err != nil {
		log.Printf("[WA] Gagal ambil data target notif: %v", err)
	} else {
		for _, s := range targetStudents {
			allTargets = append(allTargets, notifTarget{
				UserID: s.ID, FullName: s.FullName, Nisn: s.Nisn,
				ClassGroup: s.ClassGroup, ParentPhone: s.ParentPhone,
				Status: s.Status,
			})
		}
		log.Printf("[WA] Target Notif (Selain Hadir): %d siswa.", len(targetStudents))
	}

	if len(allTargets) == 0 {
		log.Println("[WA] Ga ada siswa yg perlu dinotif telat/alfa/sakit/izin.")
		return
	}

	log.Printf("[WA] Total %d siswa masuk antrian notif telat/alfa/sakit/izin.", len(allTargets))
	sent, skipped, failed := processNotificationBatch(db, allTargets, today)
	log.Printf("[WA] Done (Lainnya) — Terkirim: %d | Skip: %d | Gagal: %d", sent, skipped, failed)
}


// TestSendWhatsApp — kirim pesan test ke nomor tertentu (buat debugging)
func TestSendWhatsApp(phone, message string) (string, error) {
	if WAClient == nil || !WAClient.IsConnected() {
		return "", fmt.Errorf("WhatsApp client belum terhubung. Pastikan sudah melakukan pairing")
	}
	return SendWhatsAppMessage(phone, message)
}
