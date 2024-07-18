package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"proxyfinder/internal/collector/geonode"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

const (
	dirPath = "./storage/init/geonode"
)

var (
	verbose = true
)

// TODO: delete this shit code and write normal one
func upFillProxyTable(ctx context.Context, tx *sqlx.Tx, log *slog.Logger) error {
	// This code is executed when the migration is applied.

	_, err := os.Stat(dirPath)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	countries, err := GetCountries(ctx, tx)
	if err != nil {
		return err
	}

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

		proxies, err := geonode.ApiResponseToProxies(log, countries, content.Data)
		if err != nil {
			return err
		}

		err = SaveProxies(ctx, log, tx, proxies)
		if err != nil {
			return err
		}

	}

	return nil
}

func downFillProxyTable(ctx context.Context, tx *sqlx.Tx, log *slog.Logger) error {
	// This code is executed when the migration is rolled back.
	_, err := os.Stat(dirPath)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	countries, err := GetCountries(ctx, tx)
	if err != nil {
		return err
	}

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

		proxies, err := geonode.ApiResponseToProxies(log, countries, content.Data)
		if err != nil {
			log.Error("geonode.ApiResponseToProxies failed", slog.String("err", err.Error()))
			return err
		}

		query := "DELETE FROM proxy WHERE ip = ? AND port = ? AND protocol = ?"
		stmt, err := tx.PrepareContext(ctx, query)
		defer stmt.Close()

		for i := range proxies {
			_, err = stmt.Exec(proxies[i].Ip, proxies[i].Port, proxies[i].Protocol)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type Mute struct{}

func (m Mute) Write(p []byte) (n int, err error) {
	return 0, nil
}

func main() {
	db, err := sqlx.Open("sqlite3", "storage/local.db")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		panic(err)
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: ./initdb [up|down]")
		os.Exit(1)
	}

	log := slog.New(slog.NewTextHandler(Mute{}, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if len(args) > 2 {
		if args[2] == "--verbose" {
			log = logger.New("debug")
		}
	}

	switch args[1] {
	case "up":
		err = upFillProxyTable(context.Background(), tx, log)
	case "down":
		err = downFillProxyTable(context.Background(), tx, log)
	default:
		fmt.Println("Usage: ./initdb [up|down]")
		os.Exit(1)
	}
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	fmt.Println("Done")
}

func SaveProxies(ctx context.Context, log *slog.Logger, tx *sqlx.Tx, proxies []domain.Proxy) error {
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO proxy (ip, port, protocol, status_id, country_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range proxies {
		_, err = stmt.Exec(proxies[i].Ip, proxies[i].Port, proxies[i].Protocol, proxies[i].StatusId, proxies[i].CountryId)
		if err != nil {
			sqliteErr, ok := err.(sqlite3.Error)
			if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				log.Info("proxy already exists",
					slog.String("ip", proxies[i].Ip),
					slog.Int("port", proxies[i].Port),
					slog.String("protocol", proxies[i].Protocol))
				continue
			}
			return err
		}
	}
	return nil
}

func GetCountries(ctx context.Context, tx *sqlx.Tx) ([]domain.Country, error) {

	countries := []domain.Country{}
	res, err := tx.QueryContext(ctx, "SELECT id, name, code FROM country")
	if err != nil {
		return nil, err
	}
	for res.Next() {
		country := domain.Country{}
		err = res.Scan(&country.Id, &country.Name, &country.Code)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func GetDuplicate(proxies []domain.Proxy) []domain.Proxy {
	duplicates := []domain.Proxy{}
	for i := range proxies {
		for j := range proxies {
			if i == j {
				continue
			}
			if proxies[i].Ip == proxies[j].Ip && proxies[i].Port == proxies[j].Port && proxies[i].Protocol == proxies[j].Protocol {
				duplicates = append(duplicates, proxies[i])
			}
		}
	}
	return duplicates
}
