package user

import (
	"context"
	"errors"

	"github.com/SrJui/yap/internal/database"
	"gorm.io/gorm"
)

func Create(ctx context.Context, username string, email string, passwordHash string) error {
	user := User{Username: username, Email: email, PasswordHash: passwordHash}

	err := gorm.G[User](database.DB).Create(ctx, &user)

	return err
}

func FindByUsername(ctx context.Context, username string) (User, error) {
	if username == "" {
		return User{}, errors.New("username is required")
	}

	user := User{Username: username}

	user, err := gorm.G[User](database.DB).Where(&user).First(ctx)

	return user, err
}

func FindByEmail(ctx context.Context, email string) (User, error) {
	if email == "" {
		return User{}, errors.New("email is required")
	}

	user := User{Email: email}

	user, err := gorm.G[User](database.DB).Where(&user).First(ctx)

	return user, err
}

func UpdatePasswordHash(ctx context.Context, userID uint, passwordHash string) error {
	rows, err := gorm.G[User](database.DB).Where("id = ?", userID).Update(ctx, "password_hash", passwordHash)

	if err != nil {
		return err
	}

	if rows == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func FindByID(ctx context.Context, userID uint) (User, error) {
	return gorm.G[User](database.DB).Where("id = ?", userID).First(ctx)
}
