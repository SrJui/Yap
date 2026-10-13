package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/SrJui/yap/internal/database"
	"github.com/SrJui/yap/internal/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrInvalidResetToken = errors.New("invalid reset token")

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

func consumeResetToken(ctx context.Context, tokenHash string, passwordHash string) error {
	return database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var resetToken ResetToken
		// search data and lock transaction
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ?", tokenHash).
			First(&resetToken).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidResetToken
		}
		if err != nil {
			return err
		}
		if !resetToken.ExpiresAt.After(time.Now()) {
			return ErrInvalidResetToken
		}

		// update password hash
		result := tx.Model(&user.User{}).
			Where("id = ?", resetToken.UserID).
			Update("password_hash", passwordHash)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrInvalidResetToken
		}

		// delete all reset tokens of user
		err = tx.Where("user_id = ?", resetToken.UserID).
			Delete(&ResetToken{}).Error
		if err != nil {
			return err
		}

		// successful
		return nil
	})
}
