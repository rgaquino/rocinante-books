package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"rocinante-books/config"
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

func init() {
	if err := config.LoadJSONConfig(config.Config); err != nil {
		//logging.Fatal(logTag, "unable to load configuration. error=%v", err)
	}
	//logging.Info(logTag, "configuration file loaded")
}

func main() {
	s, err := dynamodb.New(config.Config.AWS)
	if err != nil {
		panic(err)
	}

	books, err := parseBooks(config.Config.Source.Books)
	if err != nil {
		fmt.Printf("failed to parse file, err=%v", err)
		return
	}

	booksMap := make(map[string]*entity.Book)
	for _, book := range books {
		booksMap[book.Title] = book
	}

	highlights, err := parseHighlights(config.Config.Source.Highlights)
	if err != nil {
		fmt.Printf("failed to parse file, err=%v", err)
		return
	}

	for _, highlight := range highlights {
		if book, ok := booksMap[highlight.Title]; ok {
			book.Highlights = append(book.Highlights, highlight.Quote)
		} else {
			fmt.Printf("couldn't find book: %s\n", highlight.Title)
		}
	}

	i := 0
	for _, book := range booksMap {
		if err := s.Create(book); err != nil {
			fmt.Printf("couldn't save book: %s\n, err=%v", book.Title, err)
		}
		i++
		if i == 1 {
			break
		}
	}

	f, err := os.Create("books.json")
	if err != nil {
		fmt.Printf("failed to create new file, err=%v", err)
		return
	}
	booksJson, err := json.Marshal(books)
	if _, err := f.WriteString(string(booksJson)); err != nil {
		fmt.Printf("failed to write to new file, err=%v", err)
		return
	}

	defer f.Close()
}
