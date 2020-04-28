package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-shiori/go-readability"
)

func isImage(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	mime, err := mimetype.DetectReader(resp.Body)
	if err != nil {
		return false, err
	}
	return strings.Contains(mime.String(), "image"), nil
}

func getArticle(url string) (*string, error) {
	img, err := isImage(url)
	if err != nil {
		return nil, err
	}

	if img {
		str := fmt.Sprintf("<img src=\"%s\" class=\"webfeedsFeaturedVisual \"/>", url)
		return &str, nil
	}

	article, err := readability.FromURL(url, 1*time.Second)
	if err != nil {
		return nil, err
	}
	return &article.Content, nil
}
