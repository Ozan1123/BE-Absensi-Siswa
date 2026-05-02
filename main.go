// @title My API
// @version 1.0
// @description Ini adalah dokumentasi API gue
// @host www.reihan.biz.id
// @BasePath /api/v1

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/config"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database/seeders"
	_ "github.com/KicauOrgspark/BE-Absensi-Siswa/docs" // WAJIB sesuai module
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/routes"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {

	//load env
	config.LoadEnv()

	//connect to database
	database.ConnectDB()

	database.DB.AutoMigrate(&models.Users{})

	// Inisialisasi WhatsApp client
	cfg := config.AppConfig
	waDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	if err := services.InitWA(waDSN); err != nil {
		log.Fatal("[WA] Gagal inisialisasi:", err)
	}
	if err := services.WAClient.Connect(); err != nil {
		log.Fatal("[WA] Gagal connect:", err)
	}
	log.Println("[WA] Berhasil terhubung ke server WhatsApp.")

	//to running seeders
	seeders.RunSeed()

	//start token cleaner service
	services.StartTokenCleaner()

	app := fiber.New()

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:5173,https://www.reihan.biz.id",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, //jika pake jwt
	}))

	routes.SetupRoutes(app)

	cronScheduler := services.InitAttendanceCron(database.DB)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Jalankan server di goroutine terpisah
	go func() {
		log.Printf("Server berjalan di port %s", config.AppConfig.Port)
		if err := app.Listen(config.AppConfig.Port); err != nil {
			log.Fatalf("Gagal menjalankan server: %v", err)
		}
	}()

	// Tunggu sinyal shutdown
	sig := <-quit
	log.Printf("Sinyal [%s] diterima — memulai graceful shutdown...", sig)

	// Beri waktu 10 detik untuk menyelesaikan request yang masih berjalan
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Fatalf("Gagal melakukan graceful shutdown: %v", err)
	}

	cronScheduler.Stop()
	log.Println("[CRON] Scheduler dihentikan.")

	services.WAClient.Disconnect()
	log.Println("[WA] Koneksi WhatsApp diputus.")

	log.Println("Server berhasil dimatikan dengan aman.")
	
}
