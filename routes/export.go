package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRouteExport(api fiber.Router) {
	export := api.Group("/export")

	export.Get("/attendance", handlers.ExportAttendance)
}