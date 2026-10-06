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

func FindByUsernameOrEmail(ctx context.Context, username string, email string) (User, error) {
	if username == "" && email == "" {
		return User{}, errors.New("username or email is required")
	}

	user := User{Username: username, Email: email}

	user, err := gorm.G[User](database.DB).Where(&user).First(ctx)

	return user, err
}
