package main

import (
	"log"

	"github.com/SrJui/yap/internal/auth"
	"github.com/SrJui/yap/internal/config"
	"github.com/SrJui/yap/internal/database"
	"github.com/SrJui/yap/internal/user"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	// Load config
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Get configs
	// Database config
	dbConf, err := config.LoadDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Session config
	sessConf, err := config.LoadSession()
	if err != nil {
		log.Fatal(err)
	}

	// Connect db
	database.DB, err = database.Connect(dbConf)
	if err != nil {
		log.Fatal(err)
	}

	// Migrate db
	if err = database.DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatal(err)
	}

	// Build session store
	store := auth.BuildSessionStore(dbConf, sessConf)

	app := fiber.New()

	// Use session middleware
	app.Use(auth.NewSessionMiddleware(store, sessConf))

	// Route definitions
	// Unprotected routes
	api := app.Group("api/v1")

	// Auth routes
	authGroup := api.Group("auth/")
	authGroup.Post("/signup", auth.Signup)
	authGroup.Post("/login", auth.Login)
	authGroup.Post("/logout", auth.Logout)

	// Protected routes
	protectedAPI := api.Group("", auth.RequireAuth)
	protectedAPI.Get("/test", func(c fiber.Ctx) error {
		return c.SendString("Nice!")
	})

	log.Fatal(app.Listen(":3000"))
}
