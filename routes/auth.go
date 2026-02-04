package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoute(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/login", handlers.Login)
	auth.Get("/me", middleware.ProtectedRoute, handlers.Me)
}