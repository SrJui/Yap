package user

import (
	"context"

	"github.com/SrJui/yap/internal/database"
	"gorm.io/gorm"
)

func CreateUser(username string, email string, hashedPassword string) error {
	user := User{Username: username, Email: email, PasswordHash: hashedPassword}

	ctx := context.Background()
	err := gorm.G[User](database.DB).Create(ctx, &user)

	return err
}
