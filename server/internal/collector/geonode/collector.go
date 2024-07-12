package geonode

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/storage"
	"strconv"
	"sync"
)

type Collector struct {
	url string
	log *slog.Logger
	proxyStorage storage.ProxyStorage
	countryStorage storage.CountryStorage
}

func New(log *slog.Logger, proxtStorage storage.ProxyStorage, countryStorage storage.CountryStorage) *Collector {
	return &Collector{
		url: "https://proxylist.geonode.com/api/proxy-list?protocols=http%2Chttps&limit=500&sort_by=lastChecked&sort_type=asc",
		log: log,
		proxyStorage: proxtStorage,
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

func (c *Collector) Collect(ctx context.Context) ([]domain.Proxy, error) {
	const op = "collector.geonode.Collector.Collect"

	log := c.log.With(slog.String("op", op))

	res, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	//
	// file, err := os.OpenFile("./index.json", os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// n, err := file.WriteString(string(body))
	// if err != nil {
	// 	return nil, err
	// }
	// if n < 1 {
	// 	return nil, fmt.Errorf("Error writing string to file: %d bytes", n)
	// }
	//
	// fmt.Println(string(body))
	//
	// return nil, nil

	// body, err := os.ReadFile("index.json")
	// if err != nil {
	// 	return nil, err
	// }
	
	var jsonData ApiResponse
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		log.Warn("Unmarshal error", slog.String("body", string(body)))
		return nil, err
	}
	
	proxies, err := c.ApiResponseToProxies(ctx, jsonData.Data)
	if err != nil {
		return nil, err
	}
	
	wg := sync.WaitGroup{}
	for _, v := range proxies {
		wg.Add(1)
		go func() {
			defer wg.Done()
	
			_, err := c.proxyStorage.Save(ctx, &v)
			if err != nil {
				log.Warn("ProxyStorage.Save failed", slog.Any("err", err))
			}
		}()
	}
	
	wg.Wait()


	return proxies, nil
}

func (c *Collector) ApiResponseToProxies(ctx context.Context, res []ApiProxy) ([]domain.Proxy, error) {
	countries, err := c.countryStorage.GetAll(ctx)
	if err != nil {
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
			return nil, err
		}
		proxy.Port = port
		proxy.Protocol = v.Protocols[0]
		id, ok := countryMap[v.Country]
		if !ok {
			id, err = c.countryStorage.Save(ctx, &domain.Country{Code: v.Country})
			if err != nil {
				return nil, err
			}
		}
		proxy.CountryId = id

		proxies = append(proxies, proxy)
	}

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
