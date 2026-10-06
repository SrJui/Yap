package auth

import (
	"github.com/SrJui/yap/internal/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
)

func BuildSessionStore(dbConf config.DatabaseConfig, sessConf config.SessionConfig) fiber.Storage {
	return postgres.New(postgres.Config{
		Host:     dbConf.Host,
		Username: dbConf.User,
		Password: dbConf.Password,
		Port:     dbConf.Port,
		Database: dbConf.Name,
		Table:    sessConf.Table,
		SSLMode:  dbConf.SSLMode,
	})
}

func NewSessionMiddleware(store fiber.Storage, conf config.SessionConfig) fiber.Handler {
	return session.New(session.Config{
		Storage:         store,
		CookieSecure:    conf.CookieSecure,
		CookieHTTPOnly:  true,
		CookieSameSite:  "Lax",
		IdleTimeout:     conf.IdleTimeout,
		AbsoluteTimeout: conf.AbsoluteTimeout,
		Extractor:       extractors.FromCookie("session_id"),
	})
}
