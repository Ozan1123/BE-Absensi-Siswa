package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	SetupAuthRoute(api)
	SetupRouteAttedanceToken(api)
	SetupRouteLogs(api)
	SetupRouteDashboard(api)
	SetupRouteExport(api)
}