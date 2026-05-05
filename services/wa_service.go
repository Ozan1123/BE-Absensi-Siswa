package services

import (
	"context"
	"fmt"
	"log"
	"os"
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
	if WAClient.Store.ID == nil {
		// Belum ada sesi tersimpan — pakai Pairing Code (OTP)
		phone := os.Getenv("WA_PHONE")
		if phone == "" {
			return fmt.Errorf("WA_PHONE belum diset di .env (format: 628xxxxxxxxxx)")
		}

		if err := WAClient.Connect(); err != nil {
			return fmt.Errorf("gagal connect WA: %w", err)
		}

		// Tunggu sebentar supaya koneksi stabil sebelum request pairing code
		time.Sleep(1 * time.Second)

		code, err := WAClient.PairPhone(context.Background(), phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if err != nil {
			return fmt.Errorf("gagal generate pairing code: %w", err)
		}

		log.Println("========================================")
		log.Println("[WA] PAIRING CODE (masukkan di WhatsApp)")
		log.Printf("[WA] Kode: %s", code)
		log.Println("[WA] Buka WhatsApp > Perangkat Tertaut > Tautkan Perangkat > Tautkan dengan nomor telepon")
		log.Println("========================================")

		// Tunggu sampai pairing berhasil (max 3 menit)
		timeout := time.After(3 * time.Minute)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				return fmt.Errorf("timeout menunggu pairing — coba jalankan ulang server")
			case <-ticker.C:
				if WAClient.Store.ID != nil {
					log.Println("[WA] Pairing berhasil! Sesi tersimpan.")
					return nil
				}
			}
		}
	} else {
		if err := WAClient.Connect(); err != nil {
			return fmt.Errorf("gagal reconnect WA: %w", err)
		}
		log.Println("[WA] Sesi tersimpan, berhasil terhubung kembali.")
	}
	return nil
}

// GetWAStatus — cek status koneksi WhatsApp saat ini
func GetWAStatus() map[string]interface{} {
	status := map[string]interface{}{
		"connected":    false,
		"has_session":  false,
		"phone_number": "",
	}

	if WAClient == nil {
		status["status"] = "not_initialized"
		return status
	}

	status["has_session"] = WAClient.Store.ID != nil
	status["connected"] = WAClient.IsConnected()

	if WAClient.Store.ID != nil {
		status["phone_number"] = WAClient.Store.ID.User
		status["status"] = "connected"
	} else if WAClient.IsConnected() {
		status["status"] = "waiting_pair"
	} else {
		status["status"] = "disconnected"
	}

	return status
}

// RequestPairingCode — generate pairing code baru untuk linking dari FE
func RequestPairingCode(phone string) (string, error) {
	if WAClient == nil {
		return "", fmt.Errorf("WhatsApp belum diinisialisasi")
	}

	// Kalau sudah ada sesi aktif, harus logout dulu
	if WAClient.Store.ID != nil {
		return "", fmt.Errorf("WhatsApp sudah terhubung. Logout dulu jika ingin pair ulang")
	}

	// Pastikan client connect dulu
	if !WAClient.IsConnected() {
		if err := WAClient.Connect(); err != nil {
			return "", fmt.Errorf("gagal connect WA: %w", err)
		}
		time.Sleep(1 * time.Second)
	}

	code, err := WAClient.PairPhone(context.Background(), phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		return "", fmt.Errorf("gagal generate pairing code: %w", err)
	}

	log.Printf("[WA] Pairing code baru di-generate untuk %s: %s", phone, code)

	// Monitor pairing di background
	go func() {
		timeout := time.After(3 * time.Minute)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				log.Println("[WA] Pairing code expired (3 menit).")
				return
			case <-ticker.C:
				if WAClient.Store.ID != nil {
					log.Println("[WA] Pairing berhasil via API! Sesi tersimpan.")
					return
				}
			}
		}
	}()

	return code, nil
}

// LogoutWA — disconnect dan hapus sesi WA, supaya bisa pair ulang
func LogoutWA() error {
	if WAClient == nil {
		return fmt.Errorf("WhatsApp belum diinisialisasi")
	}

	if WAClient.Store.ID == nil {
		return fmt.Errorf("tidak ada sesi WA yang aktif")
	}

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
