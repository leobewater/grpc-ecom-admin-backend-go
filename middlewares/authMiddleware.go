package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leobewater/udemy-orders-go-admin/util"
)

// IsAuthenticated gets jwt value from cookie for authentication
func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	if _, err := util.ParseJwt(cookie); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	return c.Next()
}
