package main

import (
	"encoding/json"
	"fmt"
	"os"
	"rocinante-books/config"
	"rocinante-books/data/dynamodb"
	"strconv"
)

var out string
var persist bool
var conf = &struct {
	AWS    *config.AWS    `json:"aws"`
	Source *config.Source `json:"source"`
}{}

func init() {
	var err error
	if err = config.LoadJSONConfig(conf); err != nil {
		panic(err)
	}
	persist, err = strconv.ParseBool(os.Getenv("PERSIST"))
	if err != nil {
		persist = false
	}
	out = os.Getenv("OUTPUT_FILE")
}

func main() {
	books, booksMap, err := parse(conf.Source)
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
		s, err := dynamodb.New(conf.AWS)
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
