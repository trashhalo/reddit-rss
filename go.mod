module github.com/trashhalo/reddit-rss

go 1.14

require (
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/cameronstanley/go-reddit v0.0.0-20170423222116-4bfac7ea95af
	github.com/gabriel-vasile/mimetype v1.1.1
	github.com/getsentry/sentry-go v0.6.1
	github.com/go-errors/errors v1.1.1 // indirect
	github.com/go-redis/cache v6.4.0+incompatible // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-shiori/dom v0.0.0-20200611094855-2cf8a4b8b9eb // indirect
	github.com/go-shiori/go-readability v0.0.0-20200413080041-05caea5f6592
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/gorilla/feeds v1.1.1
	github.com/graph-gophers/dataloader v5.0.0+incompatible
	github.com/jarcoal/httpmock v1.0.6 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/victorspringer/http-cache v0.0.0-20190721184638-fe78e97af707
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	gotest.tools v2.2.0+incompatible
)

replace github.com/cameronstanley/go-reddit => ./pkg/reddit
