package auth

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"github.com/SrJui/yap/internal/user"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// 3-18 chars
// _ allowed but not as first char
// . allowed but not as first char
var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9-][a-zA-Z0-9._-]{2,17}$`)

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
	if !usernamePattern.MatchString(username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username is invalid",
		})
	}

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is invalid",
		})
	}

	if !ValidatePassword(password) {
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
