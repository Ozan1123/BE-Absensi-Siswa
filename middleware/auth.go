package middleware

import (
	"strings"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware adalah middleware generic yang menerima list allowed roles.
// Contoh: AuthMiddleware("siswa", "admin", "guru", "superadmin")
func AuthMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		// Cek apakah role ada di allowedRoles
		roleAllowed := false
		for _, r := range allowedRoles {
			if r == role {
				roleAllowed = true
				break
			}
		}
		if !roleAllowed {
			return c.Status(403).JSON(fiber.Map{"error": "akses ditolak untuk role: " + role})
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
}

// Shortcut middleware untuk kemudahan pakai
var (
	// AllRoles — semua role bisa akses
	AllRoles = AuthMiddleware("siswa", "guru", "admin", "superadmin")
	// AdminOnly — hanya admin
	AdminOnly = AuthMiddleware("admin", "superadmin")
	// SiswaOnly — hanya siswa
	SiswaOnly = AuthMiddleware("siswa")
	// GuruOnly — hanya guru
	GuruOnly = AuthMiddleware("guru")
	// SuperAdminOnly — hanya superadmin
	SuperAdminOnly = AuthMiddleware("superadmin")
	// AdminGuru — admin dan guru
	AdminGuru = AuthMiddleware("admin", "guru", "superadmin")
)
