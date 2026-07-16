package handlers

import (
	"strconv"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/repo"
	"github.com/gofiber/fiber/v2"
)

func GetUnreadNotifs(c *fiber.Ctx) error {
	notifs, err := repo.GetUnreadNotifications()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil notifikasi",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil notifikasi",
		"data":    notifs,
	})
}

func ReadNotif(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
		})
	}

	err = repo.MarkAsRead(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menandai notifikasi",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notifikasi ditandai sebagai sudah dibaca",
	})
}
