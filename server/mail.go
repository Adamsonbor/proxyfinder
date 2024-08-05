package main

import (
	"fmt"
	"proxyfinder/internal/domain"
	"sync"
)

func Refresh() {
	proxies := []domain.Proxy{
		{
			Ip:       "1.1.1.1",
			Port:     80,
			Protocol: "http",
		},
		{
			Ip:       "2.2.2.2",
			Port:     80,
			Protocol: "http",
		},
		{
			Ip:       "3.3.3.3",
			Port:     80,
			Protocol: "http",
		},
	}

	wg := sync.WaitGroup{}
	for i := range proxies {
		wg.Add(1)
		go func(proxy *domain.Proxy) {
			defer wg.Done()

			proxy.Ip = "4.4.4.4"
			proxy.Port = 200
			fmt.Println(proxy)
		}(&proxies[i])
	}
	wg.Wait()

	fmt.Println()
	for _, v := range proxies {
		fmt.Println(v)
	}
}

func main() {
	Refresh()
}
