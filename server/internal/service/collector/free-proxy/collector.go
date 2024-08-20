package freeproxy

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"proxyfinder/internal/domain"
)

type Collector struct {
	url string
}

func New() *Collector {
	return &Collector{
		url: "http://free-proxy.cz/en/proxylist/country/all/http/date/all",
	}
}

func (c *Collector) Collect(ctx context.Context) ([]domain.Proxy, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Error status code: %d", res.StatusCode)
	}

	var reader io.Reader
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
		defer reader.(*gzip.Reader).Close()
	default:
		reader = res.Body
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile("index.html", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	num, err := file.WriteString(string(bodyBytes))
	if err != nil {
		return nil, err
	}
	if num < 1 {
		return nil, fmt.Errorf("Error writing to file: %d bytes", num)
	}
	return nil, nil
}
