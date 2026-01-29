package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/", handlers.Welcome)
}