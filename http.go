package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gabriel-vasile/mimetype"
	"github.com/go-shiori/go-readability"
)

type fileType int

const (
	unknown fileType = iota
	image
	video
)

func knownTypes(m *mimetype.MIME) fileType {
	if strings.HasPrefix(m.String(), "image") {
		return image
	} else if strings.HasPrefix(m.String(), "video") {
		return video
	}
	return unknown
}

func getMimeType(url string) (*mimetype.MIME, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	mime, err := mimetype.DetectReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return mime, nil
}

func cleanupUrl(url string) (string, error) {
	if strings.Contains(url, "imgur") && strings.HasSuffix(url, "gifv") {
		return strings.ReplaceAll(url, "gifv", "webm"), nil
	}

	if strings.Contains(url, "gfycat") {
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}

		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return "", err
		}

		vid, _ := doc.Find("meta[property=\"og:video\"]").Attr("content")
		return vid, nil
	}
	return url, nil
}

func getArticle(u string) (*string, error) {
	url, err := cleanupUrl(u)
	if err != nil {
		return nil, err
	}

	t, err := getMimeType(url)
	if err != nil {
		return nil, err
	}

	switch knownTypes(t) {
	case image:
		str := fmt.Sprintf("<img src=\"%s\" class=\"webfeedsFeaturedVisual \"/>", url)
		return &str, nil
	case video:
		str := fmt.Sprintf("<video><source src=\"%s\" type=\"%s\" /></video>", url, t.String())
		return &str, nil
	}

	article, err := readability.FromURL(url, 1*time.Second)
	if err != nil {
		return nil, err
	}
	return &article.Content, nil
}
