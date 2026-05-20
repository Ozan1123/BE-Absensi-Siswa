package handlers

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/gofiber/fiber/v2"
)

// GetClasses godoc
// @Summary Ambil daftar kelas
// @Description Mengambil daftar semua kelas unik dari data siswa
// @Tags classes
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /classes [get]
func GetClasses(c *fiber.Ctx) error {
	var classes []string

	err := database.DB.
		Table("users").
		Where("role = ? AND class_group IS NOT NULL AND class_group != ''", "siswa").
		Distinct("class_group").
		Order("class_group ASC").
		Pluck("class_group", &classes).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal mengambil daftar kelas"})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    classes,
	})
}
