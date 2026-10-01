package main

import (
	"log"

	"github.com/SrJui/yap/internal/auth"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	apiGroup := app.Group("api/v1")

	authGroup := apiGroup.Group("auth/")

	authGroup.Post("/signup", auth.Signup)
	authGroup.Post("/login", auth.Login)
	authGroup.Post("/logout", auth.Logout)

	log.Fatal(app.Listen(":3000"))
}
