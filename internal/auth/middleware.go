package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func RequireAuth(c fiber.Ctx) error {
	sess := session.FromContext(c)
	if sess == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Check if user is authenticated
	if sess.Get("authenticated") != true {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}
