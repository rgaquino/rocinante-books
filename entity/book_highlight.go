package entity

import "time"

// BookHighlight ...
type BookHighlight struct {
	ID        int64     `gorm:"column:id"`
	BookID    int64     `gorm:"column:book_id"`
	Content   string    `gorm:"column:content"`
	Comment   *string   `gorm:"column:comment"`
	Chapter   *string   `gorm:"column:chapter"`
	Page      *int64    `gorm:"column:page"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName ...
func (b *BookHighlight) TableName() string {
	return "book_highlight"
}
