package handler

import (
	"net/http"
	"time"

	"github.com/trashhalo/reddit-rss/pkg/client"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	httpClient := http.DefaultClient
	client.RssHandler("https://old.reddit.com", time.Now, httpClient, client.GetArticle, w, r)
}