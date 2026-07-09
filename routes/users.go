package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRouteUsers(api fiber.Router) {
	users := api.Group("/users")

	// Routes that can be accessed by both Admin and Guru (read operations)
	users.Get("/", middleware.AdminGuru, handlers.GetUsers)
	users.Get("/:id", middleware.AdminGuru, handlers.GetUserByID)

	// Routes restricted only to Admin / SuperAdmin (write/modify operations)
	users.Post("/", middleware.AdminOnly, handlers.CreateUser)
	users.Put("/:id", middleware.AdminOnly, handlers.UpdateUser)
	users.Delete("/:id", middleware.AdminOnly, handlers.DeleteUser)
	users.Post("/:id/reset-password", middleware.AdminOnly, handlers.ResetPassword)
}
