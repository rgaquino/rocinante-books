package entity

import "time"

// Book ...
type Book struct {
	ID            int64      `gorm:"column:id"`
	Title         string     `gorm:"column:title"`
	Subtitle      *string    `gorm:"column:subtitle"`
	Author        string     `gorm:"column:author"`
	Category      string     `gorm:"column:category"`
	Notes         *string    `gorm:"column:notes"`
	Slug          string     `gorm:"column:slug"`
	IsRecommended bool       `gorm:"column:is_recommended"`
	FinishedAt    time.Time  `gorm:"column:finished_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

// TableName ...
func (b *Book) TableName() string {
	return "book"
}
