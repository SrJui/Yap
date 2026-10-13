package auth

import (
	"errors"
	"log/slog"
	"net/url"
	"strings"

	"github.com/SrJui/yap/internal/config"
	"github.com/SrJui/yap/internal/mailer"
	"github.com/SrJui/yap/internal/user"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gorm.io/gorm"
)

type PasswordResetHandler struct {
	mailer *mailer.Mailer
	config config.PasswordResetConfig
}

func NewPasswordResetHandler(mailClient *mailer.Mailer, cfg config.PasswordResetConfig) *PasswordResetHandler {
	return &PasswordResetHandler{
		mailer: mailClient,
		config: cfg,
	}
}

func (h *PasswordResetHandler) ResetPassword(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "reset_password",
		"request_id", requestid.FromContext(c),
	)

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
	currentUser, err := user.FindByEmail(c.Context(), email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "If an account exists for this email, a reset link has been sent.",
		})
	}
	if err != nil {
		requestLog.Error("user lookup failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user lookup failed",
		})
	}

	// create cryptic token
	token, err := createResetToken(c.Context(), currentUser.ID)
	if err != nil {
		requestLog.Error("reset token creation failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "reset token creation failed",
		})
	}

	// send recover email
	resetURL, err := url.Parse(h.config.URL)
	if err != nil {
		requestLog.Error("creating url failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "creating url failed",
		})
	}

	query := resetURL.Query()
	query.Set("token", token)
	resetURL.RawQuery = query.Encode()
	resetLink := resetURL.String()

	body, err := mailer.RenderPasswordReset(resetLink)
	if err != nil {
		requestLog.Error("reset email rendering failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "reset email rendering failed",
		})
	}
	err = h.mailer.Send(
		[]string{currentUser.Email},
		"Reset your password",
		body,
	)
	if err != nil {
		requestLog.Error("reset email sending failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "reset email sending failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If an account exists for this email, a reset link has been sent.",
	})
}

func ChangePassword(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "change_password",
		"request_id", requestid.FromContext(c),
	)

	// get sess info
	sess := session.FromContext(c)

	// get user id from sess info
	userID, ok := sess.Get("user_id").(uint)
	if !ok || userID == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	requestLog = requestLog.With("user_id", userID)

	// get user with user id
	currentUser, err := user.FindByID(c.Context(), userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if err != nil {
		requestLog.Error("user lookup error", "err", err)
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
		requestLog.Error("password hash validation failed", "err", err)
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
		requestLog.Error("password hashing failed", "err", err)
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
		requestLog.Error("password update failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "password update failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully changed password",
	})
}
