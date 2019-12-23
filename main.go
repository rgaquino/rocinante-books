package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"rocinante-books/data/dynamodb"
	"rocinante-books/entity"
)

func parseBooks(fn string) (books []*entity.Book, err error) {
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

func main() {
	// TODO: Pass configuration
	s, err := dynamodb.New()
	if err != nil {
		panic(err)
	}

	books, err := parseBooks("/Users/rgaquino/Developer/repos/rocinante/rocinante-books/csv/books.csv")
	if err != nil {
		fmt.Printf("failed to parse file, err=%v", err)
		return
	}

	booksMap := make(map[string]*entity.Book)
	for _, book := range books {
		booksMap[book.Title] = book
	}

	highlights, err := parseHighlights("/Users/rgaquino/Developer/repos/rocinante/rocinante-books/csv/highlights.csv")
	if err != nil {
		fmt.Printf("failed to parse file, err=%v", err)
		return
	}

	for _, highlight := range highlights {
		if book, ok := booksMap[highlight.Title]; ok {
			book.Highlights = append(book.Highlights, highlight.Quote)
		} else {
			fmt.Printf("Couldn't find book: %s\n", highlight.Title)
		}
	}

	for _, book := range booksMap {
		if err := s.Create(book); err != nil {
			fmt.Printf("Couldn't save book: %s\n", book.Title)
		}
	}

	//x, _ := json.Marshal(books)
	//fmt.Println(string(x))
}
