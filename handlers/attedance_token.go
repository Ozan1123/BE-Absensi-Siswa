package handlers

import (
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/requests"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/mappers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateToken(c *fiber.Ctx) error {
	adminID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	role := c.Locals("role")
	if role != "guru" {
		return c.Status(403).JSON(fiber.Map{"error": "Hanya guru yang bisa membuat token"})
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

func CreateTokenDefault(c *fiber.Ctx) error {
	adminID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	role := c.Locals("role")
	if role != "guru" {
		return c.Status(403).JSON(fiber.Map{"error": "Hanya guru yang bisa membuat token"})
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



func SubmitToken(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	role := c.Locals("role")
	if role != "siswa" {
		return c.Status(403).JSON(fiber.Map{"error": "Hanya murid yang bisa mengisi token"})
	}

	var req requests.SubmitToken
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Muatan tidak valid"})
	}

	token, err := utils.VerifyTokenCode(req.TokenCode)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var count int64
	database.DB.Model(&models.AttedanceLogs{}).
		Where("user_id = ? AND token_id = ?", userID, token.ID).
		Count(&count)
	
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Kamu sudah  melakukan absensi",
		})
	}

	status := "hadir"
	if time.Now().After(token.LateAfter) {
		status = "telat"
	}

	log := models.AttedanceLogs{
		UserID:      userID,
		TokenID:     token.ID,
		Status:      status,
		ClockInTime: time.Now(),
	}

	if err := database.DB.Create(&log).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success To Absen",
		"status" : status,
	})
}
