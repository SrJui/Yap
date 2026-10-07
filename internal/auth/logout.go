package auth

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func Logout(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "logout",
		"request_id", requestid.FromContext(c),
	)
	sess := session.FromContext(c)

	// Complete session reset (clears all data + new session ID)
	if err := sess.Reset(); err != nil {
		requestLog.Error("session reset error", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "session error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully logged out",
	})
}
