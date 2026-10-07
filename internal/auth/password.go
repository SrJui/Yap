package auth

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/SrJui/yap/internal/user"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gorm.io/gorm"
)

func ResetPassword(c fiber.Ctx) error {
	// get email

	// check if email is in database

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
		fmt.Println("1")
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user lookup failed",
		})
	}

	// get password
	password := c.FormValue("password")
	if password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password is missing",
		})
	}

	newPassword := c.FormValue("new_password")
	if newPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "new password is missing",
		})
	}

	// validate password with user password
	isValid, _, err := argon2id.CheckHash(password, currentUser.PasswordHash)
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
	if utf8.RuneCountInString(newPassword) < 12 {
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
		fmt.Println("2")
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "password update failed",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
