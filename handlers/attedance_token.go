package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/requests"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/responses"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/mappers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/services"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
	"github.com/gofiber/fiber/v2"
)

// CreateToken godoc
// @Summary Buat token absensi (custom)
// @Description Guru/Admin membuat token dengan durasi dan toleransi keterlambatan
// @Tags token
// @Accept json
// @Produce json
// @Param request body requests.TokenReq true "Request token"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /token/create [post]
func CreateToken(c *fiber.Ctx) error {
	adminID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req requests.TokenReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Muatan tidak valid"})
	}

	token, err := utils.CreateToken(adminID, req.Duration, req.Category)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"errors": err.Error()})
	}

	// Jadwalkan notifikasi WA otomatis setelah token expired
	schedulePostTokenNotification(req.Category, req.Duration)

	return c.Status(201).JSON(fiber.Map{
		"message": "Token berhasil dibuat",
		"data":    mappers.ToTokenResponse(token),
	})
}

// CreateTokenHadir godoc
// @Summary Buat token absensi HADIR (Fase 1)
// @Description Guru/Admin membuat token khusus kehadiran tepat waktu (default 30 menit)
// @Tags token
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /token/create/hadir [post]
func CreateTokenHadir(c *fiber.Ctx) error {
	adminID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	token, err := utils.CreateToken(adminID, 30, "hadir")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"errors": err.Error()})
	}

	schedulePostTokenNotification("hadir", 30)

	return c.Status(201).JSON(fiber.Map{
		"message": "Token Hadir berhasil dibuat",
		"data":    mappers.ToTokenResponse(token),
	})
}

// CreateTokenTelat godoc
// @Summary Buat token absensi TELAT (Fase 2)
// @Description Guru/Admin membuat token khusus keterlambatan (default 60 menit)
// @Tags token
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /token/create/telat [post]
func CreateTokenTelat(c *fiber.Ctx) error {
	adminID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	token, err := utils.CreateToken(adminID, 60, "telat")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"errors": err.Error()})
	}

	schedulePostTokenNotification("telat", 60)

	return c.Status(201).JSON(fiber.Map{
		"message": "Token Telat berhasil dibuat",
		"data":    mappers.ToTokenResponse(token),
	})
}

// schedulePostTokenNotification — jadwalkan broadcast notifikasi WA setelah token expired.
// Berjalan di background goroutine, tidak memblokir response API.
func schedulePostTokenNotification(category string, durationMinutes int) {
	time.AfterFunc(time.Duration(durationMinutes)*time.Minute, func() {
		log.Printf("[WA-TIMER] Token %s expired (%d menit) — memulai broadcast notifikasi...", category, durationMinutes)
		if category == "hadir" {
			services.NotifyPresentStudents(database.DB)
		} else if category == "telat" {
			services.AutoAlfaAndNotify(database.DB)
		}
		log.Println("[WA-TIMER] Broadcast notifikasi selesai.")
	})
	log.Printf("[WA-TIMER] Notifikasi WA '%s' dijadwalkan %d menit dari sekarang.", category, durationMinutes)
}

// SubmitToken godoc
// @Summary Submit token absensi
// @Description Siswa memasukkan token untuk melakukan absensi
// @Tags token
// @Accept json
// @Produce json
// @Param request body requests.SubmitToken true "Submit token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /token/absen [post]
func SubmitToken(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req requests.SubmitToken
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Muatan tidak valid"})
	}

	// Verifikasi token
	token, isExpired, err := utils.VerifyTokenCode(req.TokenCode)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Tolak jika token sudah kedaluwarsa (harus menunggu QR fase berikutnya jika masih ada)
	if isExpired {
		return c.Status(400).JSON(fiber.Map{"error": "Token QR sudah kedaluwarsa, silakan minta QR yang baru"})
	}

	var count int64
	database.DB.Model(&models.AttedanceLogs{}).
		Where("user_id = ? AND token_id = ?", userID, token.ID).
		Count(&count)

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Kamu sudah melakukan absensi",
		})
	}

	// Tentukan status via service layer (berdasarkan kategori token)
	status := services.DetermineAttendanceStatus(token)

	now := utils.Now()
	tokenID := token.ID

	log := models.AttedanceLogs{
		UserID:      userID,
		TokenID:     &tokenID,
		Status:      status,
		ClockInTime: now,
	}

	if err := database.DB.Create(&log).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success To Absen",
		"status":  status,
	})
}

//ini untuk membuat api get qr code by id token
func GetTokenQRImage(c *fiber.Ctx) error {

	id := c.Params("id")

	var token models.AttedanceTokens

	if err := database.DB.
		First(&token, id).Error; err != nil {

		return c.Status(404).JSON(fiber.Map{
			"error": "token not found",
		})
	}

	png, err := utils.GenerateQRCode(
		token.TokenCode,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed generate qr",
		})
	}

	c.Set("Content-Type", "image/png")

	return c.Send(png)
}

//ini untuk membuat api get qr code all token aktif
func GetActiveTokens(c *fiber.Ctx) error {

	var tokens []models.AttedanceTokens

	now := time.Now()

	err := database.DB.
		Where("is_active = ?", true).
		Where("valid_until > ?", now).
		Order("created_at desc").
		Find(&tokens).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "gagal ambil token aktif",
		})
	}

	var result []fiber.Map

	for _, token := range tokens {

		result = append(result, fiber.Map{
			"id":         token.ID,
			"token_code": token.TokenCode,

			"qr_url": fmt.Sprintf(
				"/api/v1/tokens/%d/image",
				token.ID,
			),

			"expired_at": token.ValidUntil,
			"category":   token.Category,
			"is_active":  token.IsActive,
		})
	}

	return c.JSON(fiber.Map{
		"message": "success get qr code active!",
		"data":    result,
	})
}

// GetTokensPaginated godoc
// @Summary Ambil daftar semua token (paginasi)
// @Description Mengambil riwayat semua token absensi dengan paginasi
// @Tags token
// @Produce json
// @Param page query int false "Nomor halaman (default: 1)"
// @Param limit query int false "Jumlah per halaman (default: 20)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /token [get]
func GetTokensPaginated(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var total int64
	if err := database.DB.Model(&models.AttedanceTokens{}).Count(&total).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal menghitung total token"})
	}

	var tokens []models.AttedanceTokens
	if err := database.DB.
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tokens).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal mengambil data token"})
	}

	var result []responses.TokenRes
	for _, t := range tokens {
		result = append(result, mappers.ToTokenResponse(&t))
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"tokens":     result,
			"totalPages": totalPages,
			"page":       page,
			"limit":      limit,
			"total":      total,
		},
	})
}
