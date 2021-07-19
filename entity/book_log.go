package entity

import "time"

// BookLog ...
type BookLog struct {
	BookID     int64     `gorm:"column:book_id"`
	FinishedAt time.Time `gorm:"column:created_at"`
}

// TableName ...
func (b *BookLog) TableName() string {
	return "book_log"
}
