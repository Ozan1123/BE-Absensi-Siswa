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

func ReadAllNotifs(c *fiber.Ctx) error {
	err := repo.MarkAllAsRead()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menandai semua notifikasi",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Semua notifikasi ditandai sebagai sudah dibaca",
	})
}

func DeleteSelectedNotifs(c *fiber.Ctx) error {
	var payload struct {
		IDs []int64 `json:"ids"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	if len(payload.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "IDs tidak boleh kosong",
		})
	}

	err := repo.DeleteNotifications(payload.IDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus notifikasi",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notifikasi berhasil dihapus",
	})
}

func DeleteAllNotifs(c *fiber.Ctx) error {
	err := repo.DeleteAllNotifications()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus semua notifikasi",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Semua notifikasi berhasil dihapus",
	})
}

