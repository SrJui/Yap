package auth

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/SrJui/yap/internal/user"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gorm.io/gorm"
)

func Login(c fiber.Ctx) error {
	sess := session.FromContext(c)

	login := strings.TrimSpace(c.FormValue("login"))
	if login == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username or email is missing",
		})
	}

	password := c.FormValue("password")
	if password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is missing",
		})
	}
	username := ""
	email := ""
	if strings.Contains(login, "@") {
		email = login
		address, err := mail.ParseAddress(email)
		if err != nil || address.Address != email {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid email",
			})
		}
	} else {
		username = login
	}

	// find user in database
	foundUser, err := user.FindByUsernameOrEmail(c.Context(), username, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "db error",
		})
	}

	// check if password is valid
	isValid, _, err := argon2id.CheckHash(password, foundUser.PasswordHash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "hash validation failed",
		})
	}

	if !isValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}
	// set coockie or something so that user is logged in
	if err := sess.Regenerate(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error",
		})
	}

	sess.Set("user_id", foundUser.ID)
	sess.Set("authenticated", true)

	return c.SendStatus(fiber.StatusOK)
}
