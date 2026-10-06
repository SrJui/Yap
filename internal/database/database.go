package database

import (
	"fmt"

	"github.com/SrJui/yap/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(conf config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", conf.Host, conf.Port, conf.User, conf.Password, conf.Name, conf.SSLMode)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
}
