package handlers

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/services"
)

func ImportUsersExcel(c *fiber.Ctx) error {

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "file excel wajib diupload",
		})
	}

	if filepath.Ext(file.Filename) != ".xlsx" {
		return c.Status(400).JSON(fiber.Map{
			"error": "file harus berekstensi .xlsx",
		})
	}

	if file.Size == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "file kosong",
		})
	}

	tempPath := "./uploads/" + file.Filename

	if err := c.SaveFile(file, tempPath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "gagal menyimpan file",
		})
	}

	defer os.Remove(tempPath)

	result, err := services.ImportUsersFromExcel(tempPath)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Import successful",
		"data":  result,
	})
}
