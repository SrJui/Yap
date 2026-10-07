package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func Logout(c fiber.Ctx) error {
	sess := session.FromContext(c)

	// Complete session reset (clears all data + new session ID)
	if err := sess.Reset(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully logged out",
	})
}
