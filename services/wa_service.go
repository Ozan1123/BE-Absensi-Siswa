package services

import (
	"context"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var WAClient *whatsmeow.Client

// InitWA menginisialisasi WhatsApp client menggunakan whatsmeow dengan MySQL store.
// Parameter dsn menggunakan format MySQL standar: user:pass@tcp(host:port)/dbname
// CATATAN: sqlstore secara resmi hanya mendukung SQLite & PostgreSQL.
// Jika MySQL gagal, pertimbangkan untuk beralih ke SQLite.
func InitWA(dsn string) error {
	dbLog := waLog.Stdout("WA-DB", "DEBUG", true)
	waLogger := waLog.Stdout("WA-Client", "DEBUG", true)

	// Inisialisasi SQL store untuk menyimpan sesi WhatsApp
	container, err := sqlstore.New(context.Background(), "mysql", dsn, dbLog)
	if err != nil {
		log.Printf("[WA] Gagal inisialisasi store: %v", err)
		return err
	}

	// Ambil device pertama (atau buat baru jika belum ada)
	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		log.Printf("[WA] Gagal mengambil device: %v", err)
		return err
	}

	// Inisialisasi client WhatsApp
	WAClient = whatsmeow.NewClient(deviceStore, waLogger)

	log.Println("[WA] Client berhasil diinisialisasi.")
	return nil
}
