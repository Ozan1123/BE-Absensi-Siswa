package handlers

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/repo"
	"github.com/gofiber/fiber/v2"
)

func Dashboard(c *fiber.Ctx) error {
	
	data, err := repo.GetDashboardData()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error" : "failed to get dashboard data"})
	}
	return c.Status(200).JSON(fiber.Map{"Message" : "success get dashboard data", "data" : data})
}