package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	gbooks "google.golang.org/api/books/v1"
	"google.golang.org/api/option"

	"github.com/rgaquino/rocinante-books/config"
	"github.com/rgaquino/rocinante-books/entity"
)

type BookParser struct {
	imageBaseURL string
	bookService  *gbooks.Service
}

func NewBookParser(imageBaseURL string, apiKey string) *BookParser {
	opts := option.WithAPIKey(apiKey)

	bookService, err := gbooks.NewService(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	return &BookParser{
		imageBaseURL: imageBaseURL,
		bookService:  bookService,
	}
}

func (bp *BookParser) getDetails(q string, book *entity.Book) error {
	volumes, err := bp.bookService.Volumes.List(q).Do()
	if err != nil {
		return err
	}
	if len(volumes.Items) < 1 {
		return errors.New("no book details found")
	}
	details := volumes.Items[0].VolumeInfo
	book.PageCount = details.PageCount
	book.Publisher = details.Publisher
	for _, identifier := range details.IndustryIdentifiers {
		if identifier.Type == "ISBN_13" {
			book.ISBN = identifier.Identifier
		}
	}
	return nil
}

func (bp *BookParser) Parse(c *config.Source) (entity.Books, entity.BooksMap, error) {
	books, err := bp.parseBooks(c.Books)
	if err != nil {
		fmt.Printf("failed to parse file, err=%v", err)
		return nil, nil, err
	}

	booksMap := make(entity.BooksMap)
	for _, book := range books {
		booksMap[book.Title] = book
	}

	highlights, err := parseHighlights(c.Highlights)
	if err != nil {
		fmt.Printf("failed to parse file, err=%v", err)
		return nil, nil, err
	}

	for _, highlight := range highlights {
		if book, ok := booksMap[highlight.Title]; ok {
			book.Highlights = append(book.Highlights, highlight.Quote)
		} else {
			fmt.Printf("couldn't find book: %s\n", highlight.Title)
		}
	}

	return books, booksMap, nil
}

func (bp *BookParser) parseBooks(fn string) (books entity.Books, err error) {
	booksFile, err := os.Open(fn)
	if err != nil {
		return nil, nil
	}
	defer booksFile.Close()
	lines, err := csv.NewReader(booksFile).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(lines) < 1 {
		return nil, errors.New("file doesn't have contents")
	}

	for _, l := range lines[1:] {
		if l[8] == "" || !strings.Contains(l[8], "Finished") {
			continue
		}
		b := &entity.Book{
			Author:   l[2],
			Category: l[3],
		}

		fullTitle := l[0]
		fmt.Printf("processing book: %q\n", fullTitle)

		titles := strings.SplitN(fullTitle, ":", 2)
		b.Title = strings.TrimSpace(titles[0])
		if len(titles) > 1 {
			b.Subtitle = strings.TrimSpace(titles[1])
		}

		if l[6] != "" {
			finishedAt, err := time.Parse("Jan _2, 2006", l[6])
			if err != nil {
				fmt.Printf("couldn't parse finishedAt date book=%q\n", l[0])
			} else {
				b.LastFinishedAt = &finishedAt
				b.FinishedAt = []time.Time{finishedAt}
			}
		}

		// Get details
		if err := bp.getDetails(fullTitle, b); err != nil {
			fmt.Printf("failed to find details for book=%q\n, err=%v", fullTitle, err)
		}

		b.ImageLink = fmt.Sprintf("%s/%s.jpg", bp.imageBaseURL, b.ISBN)
		books = append(books, b)
	}
	return books, nil
}

type HighlightsCsv struct {
	Title  string
	Author string
	Quote  string
}

func parseHighlights(fn string) (highlights []*HighlightsCsv, err error) {
	highlightsFile, err := os.Open(fn)
	if err != nil {
		return nil, nil
	}
	defer highlightsFile.Close()
	lines, err := csv.NewReader(highlightsFile).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(lines) < 1 {
		return nil, errors.New("file doesn't have contents")
	}

	for _, l := range lines[1:] {
		h := &HighlightsCsv{
			Author: l[1],
			Quote:  l[2],
		}
		titles := strings.SplitN(l[0], ":", 2)
		h.Title = strings.TrimSpace(titles[0])
		highlights = append(highlights, h)
	}
	return highlights, nil
}
