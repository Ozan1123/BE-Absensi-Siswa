package handlers

import (
	"strconv"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/services"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/requests"
	"github.com/gofiber/fiber/v2"
)

// GetAllLogs godoc
// @Summary Ambil semua log absensi user
// @Description Mengambil riwayat absensi berdasarkan user yang sedang login
// @Tags logs
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /logs [get]
func GetAllLogs(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	query := requests.LogQuery{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search", ""),
	}

	data, err := services.GetLogService(userID, query)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(data)
}