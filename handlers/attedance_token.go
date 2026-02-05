package handlers

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/requests"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/mappers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateToken(c *fiber.Ctx) error {
	adminID := c.Locals("user_id").(int64)

	var req requests.TokenReq
	if err := c.BodyParser(&req);err != nil {
		return c.Status(400).JSON(fiber.Map{"error" : "invalid payload"})
	}

	token, err := utils.CreateToken(adminID, req.Duration)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error" : err})
	}

	return c.Status(201).JSON(fiber.Map{"Message" : "Success to create Token !", "data" : mappers.ToTokenResponse(token)})
}


func SubmitToken(c *fiber.Ctx) error {
	// userID := c.Locals("user_id").(int64)

	var req requests.SubmitToken
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error" : "invalid payload"})
	}

	if ok, _ := utils.VerifyTokenCode(req.TokenCode); !ok {
		return c.Status(400).JSON(fiber.Map{"error" : "token invalid or expired !"})
	}

	return c.Status(200).JSON(fiber.Map{"Message" : "Success To Absen"})
}