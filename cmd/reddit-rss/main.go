package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/trashhalo/reddit-rss/pkg/client"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/redis"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})

	if err != nil {
		panic(err)
	}

	log.Println("starting reddit-rss")

	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	http.HandleFunc("/info/ping", sentryHandler.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	var rssHandler http.Handler
	rssHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpClient := http.DefaultClient
		client.RssHandler("https://old.reddit.com", time.Now, httpClient, client.GetArticle, w, r)
	})

	redisCacheUrl := os.Getenv("FLY_REDIS_CACHE_URL")
	if redisCacheUrl != "" {
		u, err := url.Parse(redisCacheUrl)
		if err != nil {
			log.Fatal(err)
		}
		pass, _ := u.User.Password()
		ringOpt := &redis.RingOptions{
			Addrs: map[string]string{
				"server": u.Host,
			},
			Password: pass,
		}
		cacheClient, err := cache.NewClient(
			cache.ClientWithAdapter(redis.NewAdapter(ringOpt)),
			cache.ClientWithTTL(60*time.Minute),
			cache.ClientWithRefreshKey("opn"),
		)
		if err != nil {
			log.Fatal(err)
		}
		rssHandler = cacheClient.Middleware(rssHandler)
	}

	http.Handle("/", rssHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
