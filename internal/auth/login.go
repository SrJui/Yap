package auth

import "github.com/gofiber/fiber/v3"

func Login(c fiber.Ctx) error {
	return c.SendString("Eingeloggt!")
}
