module github.com/trashhalo/reddit-rss

go 1.14

require (
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/cameronstanley/go-reddit v0.0.0-20170423222116-4bfac7ea95af
	github.com/gabriel-vasile/mimetype v1.1.0
	github.com/getsentry/sentry-go v0.6.1
	github.com/go-shiori/go-readability v0.0.0-20200403030706-03b14e1967c5
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/gorilla/feeds v1.1.1
	github.com/kr/pretty v0.2.0 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/cameronstanley/go-reddit => ./go-reddit
