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
	"proxyfinder/internal/storage"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	ErrCountryNotFound = fmt.Errorf("country not found")
)

type Collector struct {
	Url            string
	log            *slog.Logger
	proxyStorage   storage.ProxyStorage
	countryStorage storage.CountryStorage
}

func New(log *slog.Logger, proxtStorage storage.ProxyStorage, countryStorage storage.CountryStorage) *Collector {
	log.Info("New Collector")

	return &Collector{
		Url:            "https://proxylist.geonode.com/api/proxy-list?limit=500&sort_by=lastChecked&sort_type=desc",
		log:            log,
		proxyStorage:   proxtStorage,
		countryStorage: countryStorage,
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

func (c *Collector) Collect(ctx context.Context, url string, filename string) ([]domain.Proxy, error) {
	const op = "collector.geonode.Collector.Collect"

	log := c.log.With(slog.String("op", op), slog.String("url", url))
	log.Info("Collect")

	res, err := http.Get(url)
	if err != nil {
		log.Warn("http.Get failed", slog.String("err", err.Error()))
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Warn("io.ReadAll failed", slog.String("err", err.Error()))
		return nil, err
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
		c.log.Info("PageScheduler", slog.Int("page", page))
		page++
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
func ApiResponseToProxiesX(
	ctx context.Context,
	log *slog.Logger,
	tx *sqlx.Tx,
	countryStorage storage.CountryStorage,
	res []ApiProxy) ([]domain.Proxy, error) {

	log = log.With(slog.String("op", "collector.geonode.Collector.ApiResponseToProxies"))
	log.Info("ApiResponseToProxies", slog.Int("len", len(res)))

	countries, err := countryStorage.GetAll(ctx)
	if err != nil {
		log.Error("countryStorage.GetAll failed", slog.String("err", err.Error()))
		return nil, err
	}
	countryMap := map[string]int64{}
	for _, v := range countries {
		countryMap[v.Code] = v.Id
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
		id, ok := countryMap[v.Country]
		if !ok {
			id, err = countryStorage.Savex(ctx, tx, &domain.Country{Code: v.Country})
			if err != nil {
				log.Error("countryStorage.Save failed", slog.String("err", err.Error()), slog.Any("country", v.Country))
				return nil, err
			}
		}
		proxy.CountryId = id

		proxies = append(proxies, proxy)
	}

	log.Info("ApiResponseToProxies done")

	return proxies, nil
}

func GetUniqueStrings(insts []string) []string {
	uniqe := map[string]bool{}
	out := []string{}

	for _, v := range insts {
		_, ok := uniqe[v]
		if !ok {
			uniqe[v] = true
			out = append(out, v)
		}
	}

	return out
}

func ApiResponseToProxies(
	log *slog.Logger,
	countries []domain.Country,
	res []ApiProxy,
) ([]domain.Proxy, error) {

	log.Info("ApiResponseToProxiesX", slog.Int("len", len(res)))

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

	log.Info("ApiResponseToProxiesX done")

	return proxies, nil
}
