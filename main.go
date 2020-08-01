package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/cameronstanley/go-reddit"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/feeds"
)

func linkToFeed(getArticle getArticleFn, link *reddit.Link) *feeds.Item {
	var content string
	c, _ := getArticle(link)
	if c != nil {
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

	http.HandleFunc("/", sentryHandler.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		rssHandler("https://reddit.com", time.Now, http.DefaultClient, getArticle, w, r)
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
