package post

import (
	"context"

	"github.com/SrJui/yap/internal/database"
	"gorm.io/gorm"
)

func Create(ctx context.Context, content string, authorID uint, media []Media) (Post, error) {
	createdPost := Post{AuthorID: authorID, Content: content, Media: media}

	err := gorm.G[Post](database.DB).Create(ctx, &createdPost)
	if err != nil {
		return Post{}, err
	}

	return createdPost, err
}

func DeleteSpecificPost(ctx context.Context, authorID uint, postID uint) error {
	rows, err := gorm.G[Post](database.DB).Where("id = ? AND author_id = ?", postID, authorID).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func FindByID(ctx context.Context, postID uint) (Post, error) {
	return gorm.G[Post](database.DB).
		Where("id = ?", postID).
		// preloads the media and sorts it first by position (if same pos then by id)
		Preload("Media", func(db gorm.PreloadBuilder) error {
			db.Order("position ASC").Order("id ASC")
			return nil
		}).
		First(ctx)
}

func FindAll(ctx context.Context, limit, offset int) ([]Post, error) {
	return gorm.G[Post](database.DB).
		Order("created_at DESC").Order("id DESC").
		Limit(limit).Offset(offset).
		Preload("Media", func(db gorm.PreloadBuilder) error {
			db.Order("position ASC").Order("id ASC")
			return nil
		}).Find(ctx)
}
