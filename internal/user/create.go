package user

import (
	"context"

	"github.com/SrJui/yap/internal/database"
	"gorm.io/gorm"
)

func CreateUser(username string, email string, hash string) error {
	user := User{Username: username, Email: email, PasswordHash: hash}

	ctx := context.Background()
	err := gorm.G[User](database.DB).Create(ctx, &user)

	return err
}
