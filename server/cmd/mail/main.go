package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"proxyfinder/internal/broker/rabbit"
	"proxyfinder/internal/config"
	"proxyfinder/internal/logger"
	smtpmail "proxyfinder/internal/mail/smtp"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type RabbitMessage struct {
	Email string
}

// TODO: delete this shit code and write normal one
func main() {
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	// INIT logger
	log := logger.New(cfg.Env)
	log.Info("Logger initialized")

	// connect to database
	db, err := sqlx.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// connect to mail queue
	mailQueue := rabbit.NewRabbit(cfg, "mail")
	defer mailQueue.Close()

	msgs, err := mailQueue.Consume()
	if err != nil {
		panic(err)
	}

	// Initialize mail provider
	mailProvider := smtpmail.NewMailService(&cfg.Mail)

	// consume mail queue
	log.Info("Start consuming mail queue")
	for msg := range msgs {
		var req RabbitMessage
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			log.Error("Unmarshal", slog.String("err", err.Error()))
			continue
		}

		eMsg := fmt.Sprintf("Hello from Proxpro! Your email %s is verified!", req.Email)

		// send jsondata as file to mail
		err = mailProvider.SendMail(
			context.Background(),
			req.Email,
			"Proxpro",
			[]byte(eMsg),
		)
		if err != nil {
			log.Error("SendMail", slog.String("err", err.Error()))
			continue
		}
	}
}
