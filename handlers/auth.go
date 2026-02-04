package handlers

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/requests"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)


func Login(c *fiber.Ctx) error {
	var req requests.Login
	if err := c.BodyParser(&req);err != nil {
		return c.Status(400).JSON(fiber.Map{"error" : "invalid payload"})
	}

	var user models.Users
	if err := database.DB.Where("nisn = ?", req.Nisn).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error" : "not found user with this NISN"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error" : "password invalid"})
	}

	access_token, _ := utils.GenerateJWT(user.ID, user.Role)
	

	return c.Status(200).JSON(fiber.Map{
		"Message" : "success Login",
		"access_token" : access_token,
	})




}