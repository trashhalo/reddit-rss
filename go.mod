module github.com/trashhalo/reddit-rss

go 1.14

require (
	github.com/cameronstanley/go-reddit v0.0.0-20170423222116-4bfac7ea95af
	github.com/gabriel-vasile/mimetype v1.1.0
	github.com/go-shiori/go-readability v0.0.0-20200403030706-03b14e1967c5
	github.com/gorilla/feeds v1.1.1
	github.com/kr/pretty v0.2.0 // indirect
	gocloud.dev v0.19.0
)

replace github.com/cameronstanley/go-reddit => ./go-reddit
