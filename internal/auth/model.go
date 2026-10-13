package auth

import "time"

type ResetToken struct {
	TokenHash string `gorm:"primaryKey"`
	UserID    uint
	ExpiresAt time.Time
}
