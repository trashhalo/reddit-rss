package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func NewClient(ctx context.Context, proxies []Proxy, testURL string) (*http.Client, error) {
	for _, p := range proxies {
		client, err := clientFromProxy(p)
		if err != nil {
			log.Print("cannot get client for proxy", err)
		}

		log.Println("testing proxy", p.IP, p.Port)
		ok, err := testClient(ctx, client, testURL)
		if err != nil {
			return nil, err
		}
		if ok {
			log.Println("settled on proxy", p.IP)
			return client, nil
		}
	}
	return nil, fmt.Errorf("no proxy found")
}

type Proxy struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type randomUATransport struct {
	Proxied http.RoundTripper
}

func (lrt randomUATransport) RoundTrip(req *http.Request) (res *http.Response, e error) {
	RandomUserAgent(req)
	return lrt.Proxied.RoundTrip(req)
}

func clientFromProxy(p Proxy) (*http.Client, error) {
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%s", p.IP, p.Port))
	if err != nil {
		return nil, err
	}

	t := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	return &http.Client{
		Transport: randomUATransport{t},
	}, nil
}

func testClient(ctx context.Context, client *http.Client, testURL string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", testURL, nil)
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, nil
	}
	defer resp.Body.Close()
	return resp.StatusCode < 400, nil
}
