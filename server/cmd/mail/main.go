package main

import (
	"context"
	"encoding/json"
	"fmt"

	"proxyfinder/internal/broker/rabbit"
	"proxyfinder/internal/config"
	smtpmail "proxyfinder/internal/mail/smtp"
	sqlxstorage "proxyfinder/internal/storage/v2/sqlx-storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: delete this shit code and write normal one
func main() {
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	// connect to database
	db, err := sqlx.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	storage := sqlxstorage.New(db)

	// connect to mail queue
	mailQueue := rabbit.NewRabbit(cfg, "mail")
	defer mailQueue.Close()

	msgs, err := mailQueue.Consume()
	if err != nil {
		panic(err)
	}

	// Initialize mail provider
	mailProvider := smtpmail.NewSMTPProvider(&cfg.Mail)

	// consume mail queue
	for msg := range msgs {
		proxies, err := storage.GetAll(context.Background(), nil)
		if err != nil {
			panic(err)
		}

		// convert proxies to json
		jsonData, err := json.Marshal(proxies)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(msg.Body))

		// send jsondata as file to mail
		err = mailProvider.SendMailWithAttachment(
			string(msg.Body),
			"Proxies",
			[]byte("This is a list of proxies: "),
			"proxies.json",
			jsonData,
		)
		if err != nil {
			panic(err)
		}
	}
}
