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

func Signup(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "signup",
		"request_id", requestid.FromContext(c),
	)

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
			"error": "password is invalid",
		})
	}

	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		requestLog.Error("password hashing failed", "err", err)
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
		requestLog.Error("user creation failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user creation failed",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func Login(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "login",
		"request_id", requestid.FromContext(c),
	)
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
		if !isValidEmail(email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid email",
			})
		}
	} else {
		username = login
	}

	// find user in database
	var foundUser user.User
	var err error
	if email != "" {
		foundUser, err = user.FindByEmail(c.Context(), email)
	} else {
		foundUser, err = user.FindByUsername(c.Context(), username)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}
	if err != nil {
		requestLog.Error("user lookup failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "db error",
		})
	}

	// check if password is valid
	isValid, _, err := argon2id.CheckHash(password, foundUser.PasswordHash)
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
	// set coockie or something so that user is logged in
	if err := sess.Regenerate(); err != nil {
		requestLog.Error("session regeneration error", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error",
		})
	}

	sess.Set("user_id", foundUser.ID)
	sess.Set("authenticated", true)

	return c.SendStatus(fiber.StatusOK)
}

func Logout(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "logout",
		"request_id", requestid.FromContext(c),
	)
	sess := session.FromContext(c)

	// Complete session reset (clears all data + new session ID)
	if err := sess.Reset(); err != nil {
		requestLog.Error("session reset error", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "session error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully logged out",
	})
}

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
			"error": "password is invalid",
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

func (h *PasswordResetHandler) ConfirmPasswordReset(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "confirm_password_reset",
		"request_id", requestid.FromContext(c),
	)

	// get token from form field
	token := c.FormValue("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "token is missing",
		})
	}

	// get new password from form field
	newPassword := c.FormValue("newPassword")
	if newPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "new password is missing",
		})
	}

	// validate new password
	if !isValidPassword(newPassword) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password is invalid",
		})
	}

	// hash token
	hashedToken := hashResetToken(token)

	// hash password
	hashedPassword, err := argon2id.CreateHash(newPassword, argon2id.DefaultParams)
	if err != nil {
		requestLog.Error("password hashing failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error with hashing the password",
		})
	}

	err = consumeResetToken(c.Context(), hashedToken, hashedPassword)
	if errors.Is(err, ErrInvalidResetToken) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "reset token is invalid or expired",
		})
	}
	if err != nil {
		requestLog.Error("password reset failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "password reset failed",
		})
	}

	// success message
	return c.SendStatus(fiber.StatusOK)
}
