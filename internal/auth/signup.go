package auth

import "github.com/gofiber/fiber/v3"

func Signup(c fiber.Ctx) error {
	return c.SendString("Acc erstellt!")
}
