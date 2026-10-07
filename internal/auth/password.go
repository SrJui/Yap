package auth

import (
	"errors"
	"strings"

	"github.com/SrJui/yap/internal/user"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gorm.io/gorm"
)

func ResetPassword(c fiber.Ctx) error {
	// get email
	email := strings.TrimSpace(c.FormValue("email"))
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is missing",
		})
	}

	isValid := isValidEmail(email)
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is invalid",
		})
	}

	// check if email is in database
	_, err := user.FindByEmail(c.Context(), email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// send recover email
	return c.SendStatus(fiber.StatusOK)
}

func ChangePassword(c fiber.Ctx) error {
	// get sess info
	sess := session.FromContext(c)

	// get user id from sess info
	userID, ok := sess.Get("user_id").(uint)
	if !ok || userID == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// get user with user id
	currentUser, err := user.FindByID(c.Context(), userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user lookup failed",
		})
	}

	// get password
	currentPassword := c.FormValue("currentPassword")
	if currentPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "currentPassword is missing",
		})
	}

	newPassword := c.FormValue("newPassword")
	if newPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "new password is missing",
		})
	}

	// validate password with user password
	isValid, _, err := argon2id.CheckHash(currentPassword, currentUser.PasswordHash)
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

	// validate new password
	if !isValidPassword(newPassword) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password too short",
		})
	}

	// hash password
	hashedPassword, err := argon2id.CreateHash(newPassword, argon2id.DefaultParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error with hashing the password",
		})
	}

	// change password in db
	err = user.UpdatePasswordHash(c.Context(), userID, hashedPassword)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "password update failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully changed password",
	})
}
