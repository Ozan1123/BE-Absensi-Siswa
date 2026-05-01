package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRouteLogs(api fiber.Router) {
	logs := api.Group("/logs")

	logs.Get("/", middleware.SiswaOnly, handlers.GetAllLogs)
}