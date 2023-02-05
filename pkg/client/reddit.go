package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cameronstanley/go-reddit"
	"github.com/gorilla/feeds"
	"github.com/graph-gophers/dataloader"
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

type GetArticleFn = func(client *http.Client, link *reddit.Link) (*string, error)
type NowFn = func() time.Time

func RssHandler(redditURL string, now NowFn, client *http.Client, getArticle GetArticleFn, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	resp, err := client.Do(req)
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
		Link:        &feeds.Link{Href: "https://www.reddit.com/r/rss/comments/fvg3ed/i_built_a_better_rss_feed_for_reddit/"},
		Description: "Reddit RSS feed that links directly to the content",
		Author:      &feeds.Author{Name: "Stephen Solka", Email: "stephen@solka.dev"},
		Created:     nowVal,
		Updated:     nowVal,
	}

	var limit int
	limitStr, scoreLimit := r.URL.Query()["scoreLimit"]
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

	var flair string
	flairStr, hasFlair := r.URL.Query()["flair"]
	if hasFlair {
		flair = flairStr[0]
	}

	loader := articleLoader(client, getArticle)
	var thunks []dataloader.Thunk
	for _, link := range result.Data.Children {
		if hasSafe && safe && (link.Data.Over18 || strings.ToLower(link.Data.LinkFlairText) == "nsfw") {
			continue
		}

		if scoreLimit && limit > link.Data.Score {
			continue
		}

		if hasFlair && flair != "" && link.Data.LinkFlairText != flair {
			continue
		}

		thunks = append(thunks, loader.Load(ctx, dataKey(link.Data)))
	}

	for _, thunk := range thunks {
		val, err := thunk()
		if err != nil {
			continue
		}

		item := val.(*feeds.Item)
		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	w.Header().Set("Cache-Control", "public, maxage=1800")
	io.WriteString(w, rss)
}

func linkToFeed(client *http.Client, getArticle GetArticleFn, link *reddit.Link) *feeds.Item {
	var content string
	c, _ := getArticle(client, link)
	if c != nil {
		content = *c
	}
	redditUrl := os.Getenv("REDDIT_URL")
	if redditUrl == "" {
		redditUrl = "https://old.reddit.com"
	}
	author := link.Author
	authorLink := fmt.Sprintf("%s/u/%s", redditUrl, author)
	content = fmt.Sprintf(`<p><a href="%s">u/%s</a>&nbsp;<a href="%s%s">comments</a></p> %s`, authorLink, author, redditUrl, link.Permalink, content)
	u, err := url.Parse(link.URL)
	if err == nil {
		author = u.Host
	}
	t := time.Unix(int64(link.CreatedUtc), 0)
	// if item link is to reddit, replace reddit with REDDIT_URL
	itemLink := link.URL
	if strings.HasPrefix(itemLink, "https://old.reddit.com") {
		itemLink = fmt.Sprintf(`%s%s`, redditUrl, link.Permalink)
	}
	return &feeds.Item{
		Title:   link.Title,
		Link:    &feeds.Link{Href: itemLink},
		Author:  &feeds.Author{Name: author},
		Created: t,
		Id:      link.ID,
		Content: content,
	}
}

type dataKey reddit.Link

func (k dataKey) String() string {
	l := reddit.Link(k)
	return l.ID
}

func (k dataKey) Raw() interface{} { return k }

func articleLoader(client *http.Client, getArticle GetArticleFn) *dataloader.Loader {
	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		wg := &sync.WaitGroup{}
		lock := &sync.Mutex{}
		resultMap := make(map[string]*dataloader.Result)

		for _, key := range keys {
			data := reddit.Link(key.(dataKey))
			wg.Add(1)

			go func(lock *sync.Mutex, wg *sync.WaitGroup, l reddit.Link) {
				defer wg.Done()

				item := linkToFeed(client, getArticle, &l)

				lock.Lock()
				defer lock.Unlock()
				resultMap[l.ID] = &dataloader.Result{Data: item}
			}(lock, wg, data)
		}

		wg.Wait()

		var results []*dataloader.Result
		for _, key := range keys {
			data := reddit.Link(key.(dataKey))
			results = append(results, resultMap[data.ID])
		}

		return results
	}, dataloader.WithBatchCapacity(10))
}
