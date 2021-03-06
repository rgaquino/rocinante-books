package entity

import "time"

// Book ...
type Book struct {
	ID             int64       `json:"id",dynamodbav:"_id"`
	Title          string      `json:"title",dynamodbav:"title"`
	Subtitle       string      `json:"subtitle,omitempty",dynamodbav:"subtitle,omitempty"`
	Author         string      `json:"author",dynamodbav:"author"`
	Category       string      `json:"category",dynamodbav:"category"`
	Publisher      string      `json:"publisher,omitempty",dynamodbav:"publisher,omitempty"`
	Notes          string      `json:"notes,omitempty",dynamodbav:"notes,omitempty"`
	Slug           string      `json:"slug,omitempty",dynamodbav:"slug,omitempty"`
	Highlights     []string    `json:"highlights,omitempty",dynamodbav:"highlights,omitempty"`
	FinishedAt     []time.Time `json:"finishedAt,omitempty",dynamodbav:"finishedAt,omitempty"`
	LastFinishedAt *time.Time  `json:"lastFinishedAt,omitempty",dynamodbav:"lastFinishedAt,omitempty"`
}

// TableName ...
func (b *Book) TableName() string {
	return "rocinante-books"
}

type Books []*Book

type BooksMap map[string]*Book
