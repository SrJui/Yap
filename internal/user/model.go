package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"not null;uniqueIndex"`
	Email        string `gorm:"not null;uniqueIndex"`
	PasswordHash string `gorm:"not null" json:"-"`
}
