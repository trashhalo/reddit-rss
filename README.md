# reddit-rss
[![ko-fi](https://www.ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/I3I72N2AC)  
[reddit post explaining what this is](https://www.reddit.com/r/rss/comments/fvg3ed/i_built_a_better_rss_feed_for_reddit/)


## Quick Deploy

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## installation

Your options are `docker build .` or `go build ./cmd/reddit-rss`.

## query params

-   `?safe=true` filter out nsfw posts
-   `?limit=100` filter out posts with less than 100 up votes
-   `?flair=Energy%20Products` only include posts that have that flair
