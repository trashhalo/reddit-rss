package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cameronstanley/go-reddit"
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

	return url, nil
}

func getArticle(link *reddit.Link) (*string, error) {
	u := link.URL
	// todo clean up
	if strings.Contains(u, "gfycat") {
		res, err := http.Get(u)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		img, _ := doc.Find("meta[property=\"og:image\"][content$=\".jpg\"]").Attr("content")
		vid, _ := doc.Find("meta[property=\"og:video:iframe\"]").Attr("content")
		width, _ := doc.Find("meta[property=\"og:video:width\"]").Attr("content")
		height, _ := doc.Find("meta[property=\"og:video:height\"]").Attr("content")
		str := fmt.Sprintf("<div><iframe src=\"%s\" width=\"%s\" height=\"%s\"/> <img src=\"%s\" class=\"webfeedsFeaturedVisual\"/></div>", vid, width, height, img)
		return &str, nil
	}

	if strings.Contains(u, "v.redd.it") {
		video := link.SecureMedia.RedditVideo
		if video == nil {
			return nil, errors.New("video missing from json")
		}
		str := fmt.Sprintf("<iframe src=\"%s\" width=\"%d\" height=\"%d\"/> <img src=\"%s\" class=\"webfeedsFeaturedVisual\"/>", video.FallbackURL, video.Width, video.Height, link.Thumbnail)
		return &str, nil
	}

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
