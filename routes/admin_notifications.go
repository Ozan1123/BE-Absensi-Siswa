package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupAdminNotificationRoutes(api fiber.Router) {
	api.Get("/notifications", handlers.GetUnreadNotifs)
	api.Put("/notifications/read/:id", handlers.ReadNotif)
}
