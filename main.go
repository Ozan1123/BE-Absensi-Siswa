package main

import (
	"fmt"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/config"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database/seeders"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/routes"
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
	// Setup Routes
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,https://www.reihan.biz.id", // url yang boleh akses
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",// method yang boleh dilakukan
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", //content-type header wajib
		AllowCredentials: true, //jika pake jwt
	}))
	routes.SetupRoutes(app)

	//running project
	fmt.Println("Server Is Running in Port", app.Listen(config.AppConfig.Port))
}