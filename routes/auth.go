package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoute(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/login", handlers.Login)
}