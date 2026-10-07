package auth

import (
	"errors"
	"strings"

	"github.com/SrJui/yap/internal/user"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Signup(c fiber.Ctx) error {
	username := strings.TrimSpace(c.FormValue("username"))
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username is missing",
		})
	}

	email := strings.TrimSpace(c.FormValue("email"))
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is missing",
		})
	}

	password := c.FormValue("password")
	if password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password is missing",
		})
	}

	// validation
	if !isValidUsername(username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username is invalid",
		})
	}

	if !isValidEmail(email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is invalid",
		})
	}

	if !isValidPassword(password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password too short",
		})
	}

	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error with hashing the password",
		})
	}

	err = user.Create(c.Context(), username, email, hashedPassword)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "username or email already exists",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user creation failed",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}
