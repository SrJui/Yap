package database

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return nil, errors.New("DB_HOST is missing")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		return nil, errors.New("DB_PORT is missing")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		return nil, errors.New("DB_USER is missing")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return nil, errors.New("DB_PASSWORD is missing")
	}

	name := os.Getenv("DB_NAME")
	if name == "" {
		return nil, errors.New("DB_NAME is missing")
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		return nil, errors.New("DB_SSLMODE is missing")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, name, sslmode)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
