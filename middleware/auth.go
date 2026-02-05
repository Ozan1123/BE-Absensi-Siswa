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
		return c.Status(401).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	if !strings.HasPrefix(header, "Bearer ") {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid authorization format",
		})
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")

	token, err := utils.VerifyToken(tokenStr)
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid jwt token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid jwt claims",
		})
	}

	role, ok := claims["role"].(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid role claim",
		})
	}

	if role != "siswa" {
		return c.Status(403).JSON(fiber.Map{
			"error": "kamu bukan siswa",
		})
	}

	// ✅ ambil user_id
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid user_id claim",
		})
	}

	userID := int64(userIDFloat)

	c.Locals("user_id", userID)
	c.Locals("role", role)

	return c.Next()
}

func AdminRoute(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	if header == "" {
		return c.Status(401).JSON(fiber.Map{"error": "missing authorization header"})
	}

	if !strings.HasPrefix(header, "Bearer ") {
		return c.Status(401).JSON(fiber.Map{"error": "invalid authorization format"})
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")

	token, err := utils.VerifyToken(tokenStr)
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "invalid jwt token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "invalid jwt claims"})
	}

	role, ok := claims["role"].(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "invalid role claim"})
	}

	if role != "guru" {
		return c.Status(403).JSON(fiber.Map{"error": "kamu bukan guru"})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "invalid user_id claim"})
	}

	userID := int64(userIDFloat)

	c.Locals("user_id", userID)
	c.Locals("role", role)

	return c.Next()
}
