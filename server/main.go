package main

import (
	"context"
	"encoding/json"
	"fmt"
	sqlxstorage "proxyfinder/internal/storage/v2/sqlx-storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlx.Open("sqlite3", "storage/local.db")
	if err != nil {
		panic(err)
	}

	proxyStorage := sqlxstorage.New(db)

	proxies, err := proxyStorage.GetAll(context.Background(), 1, 10)
	if err != nil {
		panic(err)
	}

	for _, proxy := range proxies {
		j, err := json.MarshalIndent(proxy, "", "\t")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(j))
	}
	
}
