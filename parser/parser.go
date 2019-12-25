package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"rocinante-books/config"
	"rocinante-books/entity"
)

func parse(c *config.Source) (entity.Books, entity.BooksMap, error) {
	books, err := parseBooks(c.Books)
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

func parseBooks(fn string) (books entity.Books, err error) {
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
		if l[6] == "" {
			continue
		}
		books = append(books, &entity.Book{
			ID:       int64(i),
			Title:    l[0],
			Author:   l[2],
			Category: l[3],
		})
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
		highlights = append(highlights, &HighlightsCsv{
			Title:  l[0],
			Author: l[1],
			Quote:  l[2],
		})
	}
	return highlights, nil
}
