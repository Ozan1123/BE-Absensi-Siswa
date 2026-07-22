package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	_ "modernc.org/sqlite"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var WAClient *whatsmeow.Client
var waContainer *sqlstore.Container

func InitWA() error {
	dbLog := waLog.Stdout("WA-DB", "DEBUG", true)
	waLogger := waLog.Stdout("WA-Client", "DEBUG", true)

	container, err := sqlstore.New(context.Background(), "sqlite", "file:whatsapp_store.db?_pragma=foreign_keys(1)", dbLog)
	if err != nil {
		return fmt.Errorf("gagal init WA store: %w", err)
	}

	waContainer = container

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return fmt.Errorf("gagal ambil device: %w", err)
	}

	WAClient = whatsmeow.NewClient(deviceStore, waLogger)
	return nil
}

func ConnectWA() error {
	if WAClient == nil {
		return fmt.Errorf("WhatsApp client belum diinisialisasi")
	}

	if WAClient.Store.ID == nil {
		log.Println("[WA] Belum ada sesi tersimpan. Siap melakukan pairing QR Code via Dashboard Admin.")
		return nil
	}

	if err := WAClient.Connect(); err != nil {
		return fmt.Errorf("gagal reconnect WA: %w", err)
	}
	log.Println("[WA] Sesi tersimpan, berhasil terhubung kembali.")
	return nil
}

var (
	currentQR      string
	currentQRMutex sync.Mutex
)

func GetCurrentQR() string {
	currentQRMutex.Lock()
	defer currentQRMutex.Unlock()
	return currentQR
}

func SetCurrentQR(qr string) {
	currentQRMutex.Lock()
	defer currentQRMutex.Unlock()
	currentQR = qr
}

// GetWAStatus — cek status koneksi WhatsApp saat ini
func GetWAStatus() map[string]interface{} {
	status := map[string]interface{}{
		"connected":    false,
		"has_session":  false,
		"phone_number": "",
		"phone":        "",
		"qr":           GetCurrentQR(),
	}

	if WAClient == nil {
		status["status"] = "not_initialized"
		return status
	}

	status["has_session"] = WAClient.Store.ID != nil
	status["connected"] = WAClient.IsConnected()

	if WAClient.Store.ID != nil {
		phone := WAClient.Store.ID.User
		status["phone_number"] = phone
		status["phone"] = phone
		status["status"] = "connected"
		SetCurrentQR("")
	} else if GetCurrentQR() != "" {
		status["status"] = "pairing"
	} else {
		status["status"] = "disconnected"
	}

	return status
}

// StartQRPairing — generate QR code baru untuk scan WhatsApp Web dari FE
func StartQRPairing() (string, error) {
	if WAClient == nil {
		return "", fmt.Errorf("WhatsApp belum diinisialisasi")
	}

	// Kalau sudah ada sesi aktif, harus logout dulu
	if WAClient.Store.ID != nil {
		return "", fmt.Errorf("WhatsApp sudah terhubung. Logout dulu jika ingin pair ulang")
	}

	if WAClient.IsConnected() {
		if qr := GetCurrentQR(); qr != "" {
			return qr, nil
		}
		WAClient.Disconnect()
	}

	qrChan, err := WAClient.GetQRChannel(context.Background())
	if err != nil {
		return "", fmt.Errorf("gagal mendapatkan QR channel: %w", err)
	}

	if err := WAClient.Connect(); err != nil {
		return "", fmt.Errorf("gagal connect WA: %w", err)
	}

	go func() {
		for item := range qrChan {
			if item.Event == "code" {
				SetCurrentQR(item.Code)
				log.Printf("[WA] QR Code diterima: %s", item.Code)
			} else if item.Event == "success" {
				SetCurrentQR("")
				log.Println("[WA] QR Pairing berhasil!")
			} else {
				log.Printf("[WA] QR Event: %s", item.Event)
			}
		}
	}()

	time.Sleep(1 * time.Second)
	qr := GetCurrentQR()
	return qr, nil
}

// RequestPairingCode — wrapper kompatibilitas ke StartQRPairing
func RequestPairingCode(phone string) (string, error) {
	return StartQRPairing()
}

// LogoutWA — disconnect dan hapus sesi WA, supaya bisa pair ulang
func LogoutWA() error {
	if WAClient == nil {
		return fmt.Errorf("WhatsApp belum diinisialisasi")
	}

	if WAClient.Store.ID == nil {
		return fmt.Errorf("tidak ada sesi WA yang aktif")
	}

	SetCurrentQR("")

	// Logout dari WhatsApp server
	if err := WAClient.Logout(context.Background()); err != nil {
		return fmt.Errorf("gagal logout WA: %w", err)
	}

	log.Println("[WA] Berhasil logout. Sesi dihapus.")

	// Re-init client dengan device store baru
	if waContainer != nil {
		deviceStore, err := waContainer.GetFirstDevice(context.Background())
		if err != nil {
			return fmt.Errorf("gagal re-init device store: %w", err)
		}
		waLogger := waLog.Stdout("WA-Client", "DEBUG", true)
		WAClient = whatsmeow.NewClient(deviceStore, waLogger)
	}

	return nil
}
