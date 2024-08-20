package httpchecker

import (
	"context"
	"log/slog"
	"fmt"
	"net/http"
	"net/url"
	"proxyfinder/internal/domain"
)


type Checker struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Checker {
	return &Checker{log: log}
}

func (c *Checker) Check(ctx context.Context, inst domain.Proxy) (bool, error) {
	// return false, nil
	const op = "checker.Check"

	proxyUrl, err := url.Parse(fmt.Sprintf("%s://%s:%d", inst.Protocol, inst.Ip, inst.Port))
	if err != nil {
		return false, err
	}
	
	log := c.log.With(slog.String("op", op))
	log.Debug(fmt.Sprintf("proxy: %s ...", proxyUrl.String()))
	
	transport := &http.Transport {
		Proxy: http.ProxyURL(proxyUrl),
	}
	
	client := &http.Client{
		Transport: transport,
	}
	
	checkUrl := "http://google.com"
	
	req, err := http.NewRequestWithContext(ctx, "GET", checkUrl, nil)
	if err != nil {
		return false, err
	}
	
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	
	log.Debug(fmt.Sprintf("proxy: %s return status code: %d", proxyUrl.String(), res.StatusCode))
	
	if res.StatusCode == http.StatusOK {
		return true, nil
	}
	
	return false, nil
}
