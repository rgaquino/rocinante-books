package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	gbooks "google.golang.org/api/books/v1"
	"google.golang.org/api/option"

	"github.com/rgaquino/rocinante-books/config"
	"github.com/rgaquino/rocinante-books/entity"
)

type BookParser struct {
	bookService *gbooks.Service
}

func NewBookParser(apiKey string) *BookParser {
	opts := option.WithAPIKey(apiKey)

	bookService, err := gbooks.NewService(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	return &BookParser{
		bookService: bookService,
	}
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

	kindleHighlights, err := parseKindleHighlights(c.Kindle)
	if err != nil {
		fmt.Printf("failed to parse file, err=%v", err)
		return nil, nil, err
	}

	for _, highlight := range kindleHighlights {
		if book, ok := booksMap[highlight.Title]; ok && highlight.Quote != "" {
			book.Highlights = append(book.Highlights, highlight.Quote)
		} else if highlight.Quote != "" {
			fmt.Printf("couldn't find book for kindle: %s\n", highlight.Title)
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

	for i, l := range lines[1:] {
		if l[9] == "" || !strings.Contains(l[9], "Finished") {
			continue
		}
		b := &entity.BookDocument{
			ID:       int64(i),
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
				b.LastFinishedAt = finishedAt
				b.FinishedAt = []time.Time{finishedAt}
			}
		}

		slug := strings.ReplaceAll(b.Title, " ", "-") + "-" + strings.ReplaceAll(b.Author, " ", "-")
		reg, err := regexp.Compile("[^a-zA-Z0-9-]+")
		if err != nil {
			log.Fatal(err)
		}
		slug = reg.ReplaceAllString(slug, "")
		b.Slug = strings.ToLower(slug)
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

func parseKindleHighlights(fn string) (highlights []*HighlightsCsv, err error) {
	highlightsFile, err := os.Open(fn)
	if err != nil {
		return nil, nil
	}
	defer highlightsFile.Close()

	var lines []string
	reader := bufio.NewReader(highlightsFile)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		lines = append(lines, string(line))
		if string(line) == "==========" {
			title := strings.Split(lines[0], "(")
			highlights = append(highlights, &HighlightsCsv{
				Title: strings.TrimSpace(title[0]),
				Quote: strings.TrimSpace(lines[3]),
			})
			lines = make([]string, 0)
		}
	}
	return highlights, nil
}
