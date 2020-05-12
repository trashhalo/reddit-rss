package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/cameronstanley/go-reddit"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/feeds"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/" {
		http.Redirect(w, r, "https://www.reddit.com/r/rss/comments/fvg3ed/i_built_a_better_rss_feed_for_reddit/", 301)
		return
	}

	log.Println(r.URL)

	url := fmt.Sprintf("https://reddit.com%s", r.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	req.Header.Add("User-Agent", "reddit-rss 1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	var result linkListing
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       fmt.Sprintf("reddit-rss %s", r.URL),
		Link:        &feeds.Link{Href: "https://github.com/trashhalo/reddit-rss"},
		Description: "Reddit RSS feed that links directly to the content",
		Author:      &feeds.Author{Name: "Stephen Solka", Email: "s@0qz.fun"},
		Created:     now,
		Updated:     now,
	}

	var limit int
	limitStr, scoreLimit := r.URL.Query()["limit"]
	if scoreLimit {
		limit, err = strconv.Atoi(limitStr[0])
		if err != nil {
			scoreLimit = false
		}
	}

	for _, link := range result.Data.Children {
		if scoreLimit && limit > link.Data.Score {
			continue
		}

		item := linkToFeed(getArticle, &link.Data)
		if err != nil {
			log.Println(err)
			continue
		}
		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	w.Header().Set("Cache-Control", "max-age=3600, public")
	io.WriteString(w, rss)
}

type getArticleFn = func(link *reddit.Link) (*string, error)

func linkToFeed(getArticle getArticleFn, link *reddit.Link) *feeds.Item {
	var content string
	c, err := getArticle(link)
	if err != nil {
		log.Println("error downloading content", err)
	} else {
		content = *c
	}
	content = fmt.Sprintf(`<p><a href="https://reddit.com%s">comments</a></p> %s`, link.Permalink, content)
	author := link.Author
	u, err := url.Parse(link.URL)
	if err == nil {
		author = u.Host
	}
	t := time.Unix(int64(link.CreatedUtc), 0)
	return &feeds.Item{
		Title:   link.Title,
		Link:    &feeds.Link{Href: link.URL},
		Author:  &feeds.Author{Name: author},
		Created: t,
		Id:      link.ID,
		Content: content,
	}
}

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})

	if err != nil {
		panic(err)
	}

	log.Println("starting reddit-rss")

	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	http.HandleFunc("/", sentryHandler.HandleFunc(handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
