package services

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var WAClient *whatsmeow.Client

func InitWA() error {
	dbLog := waLog.Stdout("WA-DB", "DEBUG", true)
	waLogger := waLog.Stdout("WA-Client", "DEBUG", true)

	container, err := sqlstore.New(context.Background(), "sqlite", "file:whatsapp_store.db?_pragma=foreign_keys(1)", dbLog)
	if err != nil {
		return fmt.Errorf("gagal init WA store: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return fmt.Errorf("gagal ambil device: %w", err)
	}

	WAClient = whatsmeow.NewClient(deviceStore, waLogger)
	return nil
}

func ConnectWA() error {
	if WAClient.Store.ID == nil {
		qrChan, _ := WAClient.GetQRChannel(context.Background())
		if err := WAClient.Connect(); err != nil {
			return fmt.Errorf("gagal connect WA: %w", err)
		}
		log.Println("[WA] Scan QR code berikut dengan WhatsApp Anda:")
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				log.Printf("[WA] QR event: %s", evt.Event)
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
