package handlers

import (
	"fmt"
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/requests"
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

	token, err := utils.CreateToken(adminID, req.Duration, req.LateAfter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"errors": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Token berhasil dibuat",
		"data":    mappers.ToTokenResponse(token),
	})
}

// CreateTokenDefault godoc
// @Summary Buat token absensi default
// @Description Guru/Admin membuat token dengan durasi default (20 menit, telat 15 menit)
// @Tags token
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /token/create/default [post]
func CreateTokenDefault(c *fiber.Ctx) error {
	adminID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	token, err := utils.CreateToken(adminID, 20, 15)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"errors": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Token berhasil dibuat",
		"data":    mappers.ToTokenResponse(token),
	})
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

	// Verifikasi token — sekarang token expired tetap dikembalikan (tidak ditolak)
	token, isExpired, err := utils.VerifyTokenCode(req.TokenCode)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
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

	// Tentukan status via service layer (clean architecture)
	status := services.DetermineAttendanceStatus(token, isExpired)

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
			"late_after":  token.LateAfter,
			"is_active":  token.IsActive,
		})
	}

	return c.JSON(fiber.Map{
		"message": "success get qr code active!",
		"data":    result,
	})
}
