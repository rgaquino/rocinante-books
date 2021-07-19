package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	config2 "github.com/rgaquino/rocinante-books/bulk/raw2json/config"
	dynamodb2 "github.com/rgaquino/rocinante-books/internal/data/dynamodb"
)

var out string
var persist bool
var conf = &struct {
	AWS         *config2.AWS         `json:"aws"`
	Source      *config2.Source      `json:"source"`
	GoogleBooks *config2.GoogleBooks `json:"google-books"`
}{}

func init() {
	var err error
	if err = config2.LoadJSONConfig(conf); err != nil || conf == nil {
		panic(err)
	}
	persist, err = strconv.ParseBool(os.Getenv("PERSIST"))
	if err != nil {
		persist = false
	}
	out = os.Getenv("OUTPUT_FILE")
}

func main() {
	bp := NewBookParser(conf.GoogleBooks.APIKey)
	books, booksMap, err := bp.Parse(conf.Source)
	if err != nil {
		panic(err)
	}
	booksJson, err := json.Marshal(books)
	if out != "" {
		f, err := os.Create(out)
		if err != nil {
			fmt.Printf("failed to create new file, err=%v", err)
			return
		}
		if _, err := f.WriteString(string(booksJson)); err != nil {
			fmt.Printf("failed to write to new file, err=%v", err)
			return
		}
		defer f.Close()
	}
	if persist {
		s, err := dynamodb2.New(conf.AWS)
		if err != nil {
			panic(err)
		}
		for _, book := range booksMap {
			if err := s.Create(book); err != nil {
				fmt.Printf("couldn't save book: %s\n, err=%v", book.Title, err)
			}
		}
	}
}
