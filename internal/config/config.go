package config

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strconv"
	"time"
)

type SessionConfig struct {
	Table           string
	IdleTimeout     time.Duration
	AbsoluteTimeout time.Duration
	CookieSecure    bool
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type MailConfig struct {
	APIKey string
	From   string
}

func LoadSession() (SessionConfig, error) {
	table := os.Getenv("SESSION_TABLE")
	if table == "" {
		return SessionConfig{}, errors.New("SESSION_TABLE is empty")
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("SESSION_IDLE_TIMEOUT"))
	if err != nil {
		return SessionConfig{}, fmt.Errorf("SESSION_IDLE_TIMEOUT: %w", err)
	}
	if idleTimeout <= 0 {
		return SessionConfig{}, errors.New("SESSION_IDLE_TIMEOUT must be positive")
	}

	absoluteTimeout, err := time.ParseDuration(os.Getenv("SESSION_ABSOLUTE_TIMEOUT"))
	if err != nil {
		return SessionConfig{}, fmt.Errorf("SESSION_ABSOLUTE_TIMEOUT: %w", err)
	}
	if absoluteTimeout <= 0 {
		return SessionConfig{}, errors.New("SESSION_ABSOLUTE_TIMEOUT must be positive")
	}
	if absoluteTimeout < idleTimeout {
		return SessionConfig{}, errors.New(
			"SESSION_ABSOLUTE_TIMEOUT must be greater than or equal to SESSION_IDLE_TIMEOUT",
		)
	}

	cookieSecure, err := strconv.ParseBool(os.Getenv("SESSION_COOKIE_SECURE"))
	if err != nil {
		return SessionConfig{}, fmt.Errorf("SESSION_COOKIE_SECURE: %w", err)
	}

	return SessionConfig{Table: table, IdleTimeout: idleTimeout, AbsoluteTimeout: absoluteTimeout, CookieSecure: cookieSecure}, nil
}

func LoadDatabase() (DatabaseConfig, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return DatabaseConfig{}, errors.New("DB_HOST is empty")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return DatabaseConfig{}, fmt.Errorf("DB_PORT: %w", err)
	}
	if port < 1 || port > 65535 {
		return DatabaseConfig{}, errors.New("DB_PORT must be between 1 and 65535")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		return DatabaseConfig{}, errors.New("DB_USER is empty")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return DatabaseConfig{}, errors.New("DB_PASSWORD is empty")
	}

	name := os.Getenv("DB_NAME")
	if name == "" {
		return DatabaseConfig{}, errors.New("DB_NAME is empty")
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		return DatabaseConfig{}, errors.New("DB_SSLMODE is empty")
	}

	return DatabaseConfig{Host: host, Port: port, User: user, Password: password, Name: name, SSLMode: sslmode}, nil
}

func LoadMail() (MailConfig, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return MailConfig{}, errors.New("RESEND_API_KEY is empty")
	}

	from := os.Getenv("MAIL_FROM")
	if from == "" {
		return MailConfig{}, errors.New("MAIL_FROM is empty")
	}
	address, err := mail.ParseAddress(from)
	if err != nil || address.Address != from {
		return MailConfig{}, errors.New("MAIL_FROM is not a valid email")
	}

	return MailConfig{APIKey: apiKey, From: from}, nil
}
