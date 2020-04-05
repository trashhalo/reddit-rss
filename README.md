reddit-rss
==========

[reddit post explaining what this is](https://www.reddit.com/r/rss/comments/fvg3ed/i_built_a_better_rss_feed_for_reddit/)

## installation
Your options are `docker build .` or `go build .`.

There is only one environment variable `CACHE_PATH` which is meant to be set to the path to where to store cache of the articles.

```
CACHE_PATH=gcs://article-cache
CACHE_PATH=file:///article-cache
```

## Todo
* I dont feel like the cache is at the right level yet. I think it would be better to put the rss output in cache so we dont hit reddit on cache hits.
* If I could get a cdn infront of this I could delete all this cache code. Im having issues getting cloudflare working with google cloud run.
