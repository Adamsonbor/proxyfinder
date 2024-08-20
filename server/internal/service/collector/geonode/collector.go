package geonode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"proxyfinder/internal/domain"
	"strconv"
	"strings"
)

var (
	ErrCountryNotFound = fmt.Errorf("country not found")
)

type Collector struct {
	Url            string
	log            *slog.Logger
}

func New(log *slog.Logger) *Collector {
	log.Info("New Collector")

	return &Collector{
		Url:            "https://proxylist.geonode.com/api/proxy-list?limit=500&sort_by=lastChecked&sort_type=desc",
		log:            log,
	}
}

type ApiResponse struct {
	Data  []ApiProxy
	Total int
	Page  int
	Limit int
}
type ApiProxy struct {
	Ip                 string `json:"ip"`
	NonymityLevel      string
	Asn                string
	City               string
	Country            string
	Created_at         string
	Google             bool
	Isp                string
	LastChecked        int
	Latency            float32
	Org                string
	Port               string
	Protocols          []string
	Speed              int
	UpTime             float64
	UpTimeSuccessCount int
	UpTimeTryCount     int
	Updated_at         string
	ResponseTime       int
}

// TODO: delete this shit code and write normal one
// Do request to geonode and save it to file by filename (filepath)
func (c *Collector) Collect(ctx context.Context, url string, filename string) ([]domain.Proxy, error) {
	const op = "collector.geonode.Collector.Collect"

	log := c.log.With(slog.String("op", op), slog.String("url", url))
	log.Info("Collect")

	res, err := http.Get(url)
	if err != nil {
		log.Warn("http.Get failed", slog.String("err", err.Error()))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Warn("http.Get failed", slog.Int("status_code", res.StatusCode))
		return nil, fmt.Errorf("http.Get failed: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Warn("io.ReadAll failed", slog.String("err", err.Error()))
		return nil, err
	}

	resp := ApiResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Warn("json.Unmarshal failed", slog.String("err", err.Error()))
		return nil, err
	}
	log.Info("json.Unmarshal done", slog.Int("len", len(resp.Data)))

	if len(resp.Data) == 0 {
		log.Warn("no proxies found")
		return nil, fmt.Errorf("no proxies found")
	}

	err = c.SaveBodyToFile(ctx, filename, body)
	if err != nil {
		log.Warn("SaveBodyToFile failed", slog.String("err", err.Error()))
		return nil, err
	}

	log.Info("Saved body to file")

	return nil, nil
}

func (c *Collector) NewPageScheduler() func() string {
	page := 0

	return func() string {
		page++
		c.log.Info("PageScheduler", slog.Int("page", page))
		return fmt.Sprintf("%s&page=%d", c.Url, page)
	}

}

func (c *Collector) SaveBodyToFile(ctx context.Context, filename string, body []byte) error {
	log := c.log.With(slog.String("op", "collector.geonode.Collector.SaveBodyToFile"))
	log.Info("SaveBodyToFile", slog.Int("len", len(body)))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("os.OpenFile failed", slog.String("err", err.Error()))
		return err
	}

	n, err := file.Write(body)
	if err != nil {
		log.Error("file.Write failed", slog.String("err", err.Error()))
		return err
	}
	if n != len(body) {
		log.Error("failed to write all bytes")
		return fmt.Errorf("failed to write all bytes")
	}

	return nil
}

func (c *Collector) SaveProxiesToFile(ctx context.Context, filename string, proxies []domain.Proxy) error {
	log := c.log.With(slog.String("op", "collector.geonode.Collector.SaveProxiesToFile"))
	log.Info("SaveProxiesToFile", slog.Int("len", len(proxies)))

	file, err := os.OpenFile("./index.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("os.OpenFile failed", slog.String("err", err.Error()))
		return err
	}

	var bytes []byte
	bytes, err = json.Marshal(proxies)
	if err != nil {
		log.Error("json.Marshal failed", slog.String("err", err.Error()))
		return err
	}

	n, err := file.Write(bytes)
	if err != nil {
		log.Error("file.Write failed", slog.String("err", err.Error()))
		return err
	}
	if n != len(bytes) {
		log.Error("failed to write all bytes")
		return fmt.Errorf("failed to write all bytes")
	}

	return nil
}

// TODO: delete this shit code and write normal one
func ApiResponseToProxies(
	log *slog.Logger,
	countries []domain.Country,
	res []ApiProxy,
) ([]domain.Proxy, error) {

	log.Info("ApiResponseToProxies", slog.Int("len", len(res)))

	countryMap := map[string]int64{}
	for _, v := range countries {
		countryMap[strings.ToLower(v.Code)] = v.Id
	}

	proxies := []domain.Proxy{}
	for _, v := range res {
		proxy := domain.Proxy{}
		proxy.Ip = v.Ip
		port, err := strconv.Atoi(v.Port)
		if err != nil {
			log.Error("strconv.Atoi failed", slog.String("err", err.Error()))
			return nil, err
		}
		proxy.Port = port
		proxy.Protocol = v.Protocols[0]
		id, ok := countryMap[strings.ToLower(v.Country)]
		if !ok {
			log.Error("country not found", slog.String("country", v.Country))
			return nil, ErrCountryNotFound
		}
		proxy.CountryId = id
		proxy.StatusId = domain.STATUS_UNAVAILABLE

		proxies = append(proxies, proxy)
	}

	log.Info("ApiResponseToProxies done")

	return proxies, nil
}
