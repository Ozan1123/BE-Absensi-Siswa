package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupAdminNotificationRoutes(api fiber.Router) {
	api.Get("/notifications", handlers.GetUnreadNotifs)
	api.Put("/notifications/read/:id", handlers.ReadNotif)
	api.Put("/notifications/read-all", handlers.ReadAllNotifs)
	api.Delete("/notifications/bulk", handlers.DeleteSelectedNotifs)
	api.Delete("/notifications/all", handlers.DeleteAllNotifs)
}
