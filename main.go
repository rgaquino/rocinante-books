package main

import (
	"rocinante-books/data/dynamodb"
)

func main() {
	// TODO: Pass configuration
	_, err := dynamodb.New()
	if err != nil {
		panic(err)
	}

	// TODO: Parse books.csv
	// TODO: Parse highlights.csv and add to books.csv
}
