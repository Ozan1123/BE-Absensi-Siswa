package routes

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/handlers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRouteAttedanceToken(api fiber.Router) {
	token := api.Group("/token")

	token.Post("/create", middleware.AdminOnly, handlers.CreateToken)
	token.Post("/create/default", middleware.AdminOnly, handlers.CreateTokenDefault)
	token.Post("/absen", middleware.SiswaOnly, handlers.SubmitToken)

	token.Get("/qr_code/active", handlers.GetActiveTokens)
	token.Get("/:id/image", handlers.GetTokenQRImage)
}