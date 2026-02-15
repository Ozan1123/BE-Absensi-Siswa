package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRouteDashboard(api fiber.Router) {
	api.Get("/dashboard", middleware.AdminRoute ,handlers.Dashboard)
}