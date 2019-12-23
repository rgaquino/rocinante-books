package entity

import "time"

// Book ...
type Book struct {
	ID         int64        `json:"id"`
	Title      string       `json:"title"`
	Author     string       `json:"author"`
	Category   *string      `json:"category"`
	Highlights *[]string    `json:"highlights"`
	ReadDates  *[]time.Time `json:"read_dates"`
}

// TableName ...
func (b *Book) TableName() string {
	return "rocinante-books"
}
