package main

import (
	"log"

	"github.com/SrJui/yap/internal/auth"
	"github.com/SrJui/yap/internal/database"
	"github.com/SrJui/yap/internal/user"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	database.DB, err = database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	err = database.DB.AutoMigrate(
		&user.User{},
	)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	apiGroup := app.Group("api/v1")

	authGroup := apiGroup.Group("auth/")

	authGroup.Post("/signup", auth.Signup)
	authGroup.Post("/login", auth.Login)
	authGroup.Post("/logout", auth.Logout)

	log.Fatal(app.Listen(":3000"))
}
