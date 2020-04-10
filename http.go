package main

import (
	"time"

	"github.com/go-shiori/go-readability"
)

func getArticle(url string) (*string, error) {
	article, err := readability.FromURL(url, 1*time.Second)
	if err != nil {
		return nil, err
	}
	return &article.Content, nil
}
