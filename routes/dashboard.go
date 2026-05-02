package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRouteDashboard(api fiber.Router) {
	api.Get("/dashboard", middleware.AdminOnly, handlers.Dashboard)
	api.Get("/dashboard/trend", middleware.AdminOnly, handlers.GetTrendAttendance)
}