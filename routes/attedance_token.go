package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRouteAttedanceToken(api fiber.Router) {
	token := api.Group("/token")

	token.Post("/create", middleware.AdminRoute, handlers.CreateToken)
	token.Post("/create/default", middleware.AdminRoute, handlers.CreateTokenDefault)
	token.Post("/absen", middleware.ProtectedRoute, handlers.SubmitToken)
}