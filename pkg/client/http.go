package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	gReddit "github.com/cameronstanley/go-reddit"
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

func getMimeType(client *http.Client, url string) (*mimetype.MIME, error) {
	resp, err := client.Get(url)
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

func fixAmp(url string) string {
	return strings.Replace(url, "&amp;", "&", -1)
}

func imgElement(media gReddit.MediaMetadata) string {
	if media.S.Gif != "" {
		return fmt.Sprintf("<img src=\"%s\" /><br/>", fixAmp(media.S.Gif))
	} else if media.S.U != "" {
		return fmt.Sprintf("<img src=\"%s\" /><br/>", fixAmp(media.S.U))
	} else {
		return ""
	}
}

var videoMissingErr = errors.New("video missing from json")

func GetArticle(client *http.Client, link *gReddit.Link) (*string, error) {
	u := link.URL

	if link.IsSelf {
		str := html.UnescapeString(link.SelftextHTML)
		return &str, nil
	}

	if link.Media.Oembed.Type == "video" && link.Media.Oembed.HTML != "" {
		str := html.UnescapeString(link.Media.Oembed.HTML)
		return &str, nil
	}

	if len(link.MediaMetadata) > 0 {
		var b strings.Builder
		b.WriteString("<div>")
		if len(link.GalleryData.Items) > 0 {
			for _, item := range link.GalleryData.Items {
				b.WriteString(imgElement(link.MediaMetadata[item.MediaID]))
			}
		} else {
			for _, media := range link.MediaMetadata {
				b.WriteString(imgElement(media))
			}
		}
		b.WriteString("</div>")
		str := b.String()
		return &str, nil
	}

	// todo clean up
	if strings.Contains(u, "gfycat") {
		res, err := client.Get(u)
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
			if len(link.CrossPostParentList) == 0 {
				return nil, videoMissingErr
			}
			parent := link.CrossPostParentList[0]
			video = parent.SecureMedia.RedditVideo
		}
		if video == nil {
			return nil, videoMissingErr
		}
		str := fmt.Sprintf("<iframe src=\"%s\" width=\"%d\" height=\"%d\"/> <img src=\"%s\" class=\"webfeedsFeaturedVisual\"/>", video.FallbackURL, video.Width, video.Height, link.Thumbnail)
		return &str, nil
	}

	url, err := cleanupUrl(u)
	if err != nil {
		return nil, err
	}

	t, err := getMimeType(client, url)
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

func articleFromURL(ctx context.Context, client *http.Client, pageURL string, timeout time.Duration) (readability.Article, error) {
	// Make sure URL is valid
	_, err := url.ParseRequestURI(pageURL)
	if err != nil {
		return readability.Article{}, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Fetch page from URL
	req, err := http.NewRequestWithContext(ctx, "GET", pageURL, nil)
	if err != nil {
		return readability.Article{}, fmt.Errorf("failed to create req for page: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return readability.Article{}, fmt.Errorf("failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	// Make sure content type is HTML
	cp := resp.Header.Get("Content-Type")
	if !strings.Contains(cp, "text/html") {
		return readability.Article{}, fmt.Errorf("URL is not a HTML document")
	}

	// Check if the page is readable
	var buffer bytes.Buffer
	tee := io.TeeReader(resp.Body, &buffer)

	parser := readability.NewParser()
	if !parser.IsReadable(tee) {
		return readability.Article{}, fmt.Errorf("the page is not readable")
	}

	// Parse content
	return parser.Parse(&buffer, pageURL)
}
