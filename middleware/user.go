package middleware

import (
	"strings"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ProtectedRoute(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	if header == "" {
		return c.Status(401).JSON(fiber.Map{"error": "missing header"})
	}

	if !strings.HasPrefix(header, "Bearer ") {
		return c.Status(401).JSON(fiber.Map{"error": "invalid fomating header"})
	}

	tokenStr := strings.TrimPrefix(header, "Bearer")

	token, err := utils.VerifyToken(tokenStr)
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "invalid jwt"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "invalid claims"})
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "invalid claims"})
	}

	c.Locals("user_id", int64(userID))
	return c.Next()
}
