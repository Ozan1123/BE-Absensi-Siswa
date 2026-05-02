package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupNotificationRoutes(api fiber.Router) {
	notif := api.Group("/notification")

	// Settings CRUD — hanya admin/superadmin
	notif.Get("/settings", middleware.AdminOnly, handlers.GetNotificationSettings)
	notif.Put("/settings", middleware.AdminOnly, handlers.UpdateNotificationSettings)

	// Test kirim WA — hanya admin/superadmin
	notif.Post("/test", middleware.AdminOnly, handlers.TestSendWA)

	// Log notifikasi — hanya admin/superadmin
	notif.Get("/logs", middleware.AdminOnly, handlers.GetNotificationLogs)

	// Trigger manual — hanya admin/superadmin
	notif.Post("/trigger", middleware.AdminOnly, handlers.TriggerNotificationNow)

	// Set status absensi siswa — guru/admin/superadmin
	attendance := api.Group("/attendance")
	attendance.Get("/students", middleware.AdminGuru, handlers.GetStudentsAttendanceToday)
	attendance.Put("/status", middleware.AdminGuru, handlers.UpdateStudentStatus)
}
