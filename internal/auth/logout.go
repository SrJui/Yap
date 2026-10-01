package auth

import "github.com/gofiber/fiber/v3"

func Logout(c fiber.Ctx) error {
	return c.SendString("Ausgeloggt!")
}
