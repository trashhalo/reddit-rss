package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cameronstanley/go-reddit"
	"github.com/gorilla/feeds"
)

type linkListingChildren struct {
	Kind string      `json:"kind"`
	Data reddit.Link `json:"data"`
}

type linkListingData struct {
	Modhash  string                `json:"modhash"`
	Children []linkListingChildren `json:"children"`
	After    string                `json:"after"`
	Before   interface{}           `json:"before"`
}

type linkListing struct {
	Kind string          `json:"kind"`
	Data linkListingData `json:"data"`
}

type getArticleFn = func(link *reddit.Link) (*string, error)
type nowFn = func() time.Time

func rssHandler(redditURL string, now nowFn, client *http.Client, getArticle getArticleFn, w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/" {
		http.Redirect(w, r, "https://www.reddit.com/r/rss/comments/fvg3ed/i_built_a_better_rss_feed_for_reddit/", 301)
		return
	}

	log.Println(r.URL)

	url := fmt.Sprintf("%s%s", redditURL, r.URL)
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

	nowVal := now()
	feed := &feeds.Feed{
		Title:       fmt.Sprintf("reddit-rss %s", r.URL),
		Link:        &feeds.Link{Href: "https://github.com/trashhalo/reddit-rss"},
		Description: "Reddit RSS feed that links directly to the content",
		Author:      &feeds.Author{Name: "Stephen Solka", Email: "s@0qz.fun"},
		Created:     nowVal,
		Updated:     nowVal,
	}

	var limit int
	limitStr, scoreLimit := r.URL.Query()["limit"]
	if scoreLimit {
		limit, err = strconv.Atoi(limitStr[0])
		if err != nil {
			scoreLimit = false
		}
	}

	var safe bool
	safeStr, hasSafe := r.URL.Query()["safe"]
	if hasSafe {
		safe = strings.ToLower(safeStr[0]) == "true"
	}

	for _, link := range result.Data.Children {
		if hasSafe && safe && (link.Data.Over18 || strings.ToLower(link.Data.LinkFlairText) == "nsfw") {
			continue
		}

		if scoreLimit && limit > link.Data.Score {
			continue
		}

		item := linkToFeed(getArticle, &link.Data)
		if err != nil {
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
