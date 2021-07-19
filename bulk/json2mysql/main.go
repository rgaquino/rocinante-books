package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/rgaquino/rocinante-books/entity"
	mysql2 "github.com/rgaquino/rocinante-books/internal/data/mysql"
	ptr2 "github.com/rgaquino/rocinante-books/internal/ptr"
)

func main() {
	f, err := os.Open("/Users/rgaquino/Developer/repos/rocinante/rocinante-books/migrator/export.json")
	if err != nil {
		panic(err)
	}
	booksExport, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	bookDocuments := make([]*entity.BookDocument, 0)
	if err := json.Unmarshal(booksExport, &bookDocuments); err != nil {
		panic(err)
	}

	s, err := mysql2.New("DSN")
	if err != nil {
		panic(err)
	}

	sort.SliceStable(bookDocuments, func(i, j int) bool {
		return bookDocuments[i].LastFinishedAt.Before(bookDocuments[j].LastFinishedAt)
	})

	for _, bd := range bookDocuments {
		fmt.Println(bd)

		b := &entity.Book{
			Title:         bd.Title,
			Subtitle:      ptr2.StrRefDefaultNil(bd.Subtitle),
			Author:        bd.Author,
			Category:      bd.Category,
			Notes:         bd.Notes,
			Slug:          bd.Slug,
			IsRecommended: false,
			FinishedAt:    bd.LastFinishedAt,
		}
		if err := s.Create(b); err != nil {
			fmt.Println(err)
		}

		for _, h := range bd.Highlights {
			bh := &entity.BookHighlight{
				BookID:  b.ID,
				Content: h,
			}
			if err := s.Create(bh); err != nil {
				fmt.Println(err)
			}
		}

		if len(bd.FinishedAt) < 2 {
			continue
		}

		for _, f := range bd.FinishedAt[1:] {
			bl := &entity.BookLog{
				BookID:     b.ID,
				FinishedAt: f,
			}
			if err := s.Create(bl); err != nil {
				fmt.Println(err)
			}
		}
	}
}
