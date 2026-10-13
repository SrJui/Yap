package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/SrJui/yap/internal/database"
)

func generateResetToken() string {
	tokenNumber := make([]byte, 32)
	rand.Read(tokenNumber)
	return base64.RawURLEncoding.EncodeToString(tokenNumber)
}

func hashResetToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func createResetToken(ctx context.Context, userID uint) (string, error) {
	token := generateResetToken()
	resetToken := ResetToken{TokenHash: hashResetToken(token), UserID: userID, ExpiresAt: time.Now().Add(15 * time.Minute)}
	err := database.DB.WithContext(ctx).Create(&resetToken).Error
	if err != nil {
		return "", err
	}

	return token, nil
}
