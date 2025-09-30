package middleware

import (
	"arek-muhammadiyah-be/helper/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Authorization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Authorization header required",
			})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role_id", claims.RoleID)
		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID := c.Locals("role_id").(*uint)
		if roleID == nil || *roleID != 1 { // Assuming 1 is admin role
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Admin access required",
			})
		}
		return c.Next()
	}
}