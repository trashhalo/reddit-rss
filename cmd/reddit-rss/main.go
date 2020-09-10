package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/trashhalo/reddit-rss/pkg/proxy"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/trashhalo/reddit-rss/pkg/client"
)

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
		ctx := r.Context()
		httpClient := http.DefaultClient
		proxies := proxy.ProxyList[:]
		h, err := proxy.NewClient(ctx, proxies, "https://ifconfig.io")
		if err != nil {
			log.Println("failed to get proxy client", err)
		} else {
			httpClient = h
		}

		client.RssHandler("https://reddit.com", time.Now, httpClient, client.GetArticle, w, r)
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
