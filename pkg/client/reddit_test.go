package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cameronstanley/go-reddit"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func Test_RssHandler(t *testing.T) {
	type args struct {
		getArticle GetArticleFn
		r          *http.Request
	}
	tests := []struct {
		name     string
		args     args
		reddit   linkListing
		response string
	}{
		{
			name: "Happy Path",
			args: args{
				r: httptest.NewRequest("get", "/r/android.json", nil),
				getArticle: func(client *http.Client, link *reddit.Link) (*string, error) {
					resp := "foo"
					return &resp, nil
				},
			},
			reddit: linkListing{
				Data: linkListingData{
					Children: []linkListingChildren{
						{
							Data: reddit.Link{
								Title: "My cool link",
							},
						},
					},
				},
			},
			response: "fixtures/android.xml",
		},
		{
			name: "Query Param Safe",
			args: args{
				r: httptest.NewRequest("get", "/r/android.json?safe=true", nil),
				getArticle: func(client *http.Client, link *reddit.Link) (*string, error) {
					resp := "foo"
					return &resp, nil
				},
			},
			reddit: linkListing{
				Data: linkListingData{
					Children: []linkListingChildren{
						{
							Data: reddit.Link{
								Title:  "nsfw1",
								Over18: true,
							},
						},
						{
							Data: reddit.Link{
								Title:  "sfw",
								Over18: false,
							},
						},
						{
							Data: reddit.Link{
								Title:         "nsfw2",
								LinkFlairText: "nsfw",
							},
						},
					},
				},
			},
			response: "fixtures/nsfw.xml",
		},
		{
			name: "Query Param Limit",
			args: args{
				r: httptest.NewRequest("get", "/r/android.json?limit=100", nil),
				getArticle: func(client *http.Client, link *reddit.Link) (*string, error) {
					resp := "foo"
					return &resp, nil
				},
			},
			reddit: linkListing{
				Data: linkListingData{
					Children: []linkListingChildren{
						{
							Data: reddit.Link{
								Title: "low1",
								Score: 10,
							},
						},
						{
							Data: reddit.Link{
								Title: "high",
								Score: 101,
							},
						},
						{
							Data: reddit.Link{
								Title: "low2",
								Score: 20,
							},
						},
					},
				},
			},
			response: "fixtures/score.xml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				enc := json.NewEncoder(w)
				err := enc.Encode(&tt.reddit)
				require.NoError(t, err)
			}))
			defer server.Close()

			nowVal, _ := time.Parse(time.RFC822, time.RFC822)
			now := func() time.Time {
				return nowVal
			}

			RssHandler(server.URL, now, server.Client(), tt.args.getArticle, w, tt.args.r)
			res := w.Result()
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			if os.Getenv("UPDATE_GOLDEN_FILES") == "true" {
				err := ioutil.WriteFile(tt.response, body, 0666)
				require.NoError(t, err)
			}

			fixture, err := ioutil.ReadFile(tt.response)
			require.NoError(t, err)

			assert.Equal(t, string(fixture), string(body))
		})
	}
}
