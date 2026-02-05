package handlers

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/mappers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllLogs(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	var logs []models.AttedanceLogs
	if err := database.DB.Preload("User").Where("user_id = ?", userID).Find(&logs).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error" : "not found user logs"})
	}

	return c.Status(200).JSON(fiber.Map{"Message" : "Found User Logs", "data" : mappers.ListToLogsResponse(logs)})
}