package main

import (
	"fmt"
	"os"
	"proxyfinder/internal/config"
	"proxyfinder/internal/mail/smtp"
)

func main() {
	cfg := config.MustLoadConfig()
	fmt.Println(cfg.Email)

	gmailSender := smtpmail.NewSMTPProvider(&cfg.Email)

	subject := "Email from proxyfinder"
	body := "This is a test email from proxyfinder"

	err := gmailSender.SendMail(os.Args[len(os.Args)-1], subject, body)
	if err != nil {
		panic(err)
	}
}
