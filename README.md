# reddit-rss
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/trashhalo/reddit-rss)
[![ko-fi](https://www.ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/I3I72N2AC)  
[reddit post explaining what this is](https://www.reddit.com/r/rss/comments/fvg3ed/i_built_a_better_rss_feed_for_reddit/)


## installation

Your options are `docker build .` or `go build ./cmd/reddit-rss`.

## using my free hosted version

I run a version of reddit-rss at https://reddit.0qz.fun

If you are interested in using it to you:
1. Go to a subreddit or meta feed you like example: https://www.reddit.com/r/Android/
2. Add .json onto the end: https://www.reddit.com/r/Android.json
3. Change the domain name to, reddit.0qz.fun like: https://reddit.0qz.fun/r/android.json
4. Subscribe to ^^^ that url in your favorite feed reader.

## query params

-   `?safe=true` filter out nsfw posts
-   `?limit=100` filter out posts with less than 100 up votes
-   `?flair=Energy%20Products` only include posts that have that flair

## Quick Deploy

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)
