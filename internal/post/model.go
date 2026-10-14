package post

import "time"

type Post struct {
	ID        uint      `json:"id"`
	AuthorID  uint      `gorm:"not null;index" json:"author_id"`
	Content   string    `gorm:"not null" json:"content"`
	Media     []Media   `json:"media"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Media struct {
	ID          uint   `json:"id"`
	PostID      uint   `gorm:"not null;index" json:"post_id"`
	StorageKey  string `gorm:"not null" json:"-"`
	ContentType string `gorm:"not null" json:"content_type"`
	Position    uint   `gorm:"not null" json:"position"`
}
