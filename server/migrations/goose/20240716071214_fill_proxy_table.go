package migrations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"proxyfinder/internal/collector/geonode"
	"proxyfinder/internal/storage/sqlite-storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

const (
	dirPath = "./storage/init"
)

func init() {
	goose.AddMigrationContext(upFillProxyTable, downFillProxyTable)
}

// TODO: delete this shit code and write normal one
func upFillProxyTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.

	_, err := os.Stat(dirPath)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	countryStorage := sqlite.NewCountry(nil)
	proxyStorage := sqlite.NewProxy(nil)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name()[0] == '.' {
			continue
		}
		if file.Name()[:4] != "prox" {
			continue
		}

		filePath := dirPath + "/" + file.Name()
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		var content geonode.ApiResponse
		err = json.Unmarshal(fileContent, &content)
		if err != nil {
			return err
		}

		if len(content.Data) == 0 {
			return fmt.Errorf("len(content.Data) == 0")
		}

		proxies, err := geonode.ApiResponseToProxies(ctx, nil, tx, countryStorage, content.Data)
		if err != nil {
			return err
		}

		err = proxyStorage.SaveAllx(ctx, tx, proxies)
		if err != nil {
			return err
		}
	}

	return nil
}

func downFillProxyTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := os.Stat(dirPath)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	countryStorage := sqlite.NewCountry(nil)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name()[0] == '.' {
			continue
		}
		if file.Name()[:4] != "prox" {
			continue
		}

		filePath := dirPath + "/" + file.Name()
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		var content geonode.ApiResponse
		err = json.Unmarshal(fileContent, &content)
		if err != nil {
			return err
		}

		if len(content.Data) == 0 {
			return fmt.Errorf("len(content.Data) == 0")
		}

		proxies, err := geonode.ApiResponseToProxies(ctx, nil, tx, countryStorage, content.Data)
		if err != nil {
			return err
		}

		query := "DELETE FROM proxy WHERE ip = ? AND port = ?"
		stmt, err := tx.PrepareContext(ctx, query)
		defer stmt.Close()

		for i := range proxies {
			_, err = stmt.Exec(proxies[i].Ip, proxies[i].Port)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
