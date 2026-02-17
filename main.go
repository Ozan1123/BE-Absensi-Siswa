package main

import (
	"fmt"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/config"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database/seeders"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/routes"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	//load env
	config.LoadEnv()

	//connect to database
	database.ConnectDB()

	//to running seeders
	seeders.RunSeed()

	//start token cleaner service
	services.StartTokenCleaner()

	// Setup Routes
	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:3000,https://www.reihan.biz.id",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, //jika pake jwt
	}))

	routes.SetupRoutes(app)

	//running project
	fmt.Println("Server Is Running in Port", app.Listen(config.AppConfig.Port))
}
