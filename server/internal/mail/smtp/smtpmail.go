package smtpmail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"proxyfinder/internal/config"
)

type EmailProvider interface {
	SendMail(to string, subject string, body string) error
	SendMails(to []string, subject string, body string) error
}

type SMTPProvider struct {
	Config *config.Email
	Auth   smtp.Auth
}

func NewSMTPProvider(cfg *config.Email) *SMTPProvider {
	return &SMTPProvider{
		Config: cfg,
		Auth:   smtp.PlainAuth("", cfg.Email, cfg.Pass, cfg.Addr),
	}
}

func (s *SMTPProvider) SendMail(to string, subject string, body string) error {
	return s.SendMails([]string{to}, subject, body)
}

func (s *SMTPProvider) SendMails(to []string, subject string, body string) error {
	message := bytes.Buffer{}
	message.WriteString(fmt.Sprintf("From: %s\n", s.Config.From))
	message.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	message.WriteString(body)

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.Config.Addr, s.Config.Port),
		s.Auth,
		s.Config.Email,
		to,
		message.Bytes(),
	)
}
