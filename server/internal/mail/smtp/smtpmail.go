package smtpmail

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/smtp"
	"proxyfinder/internal/config"
	"strings"
)

type Mailer interface {
	SendMail(ctx context.Context, to string, subject string, body string) error
	SendMails(ctx context.Context, to []string, subject string, body string) error
	SendMailWithAttachment(to string, subject string, body string, filename string, attachmentData []byte) error
	SendMailsWithAttachment(to []string, subject string, body string, filename string, attachmentData []byte) error
}

type MailService struct {
	log    *slog.Logger
	Config *config.Mail
	Auth   smtp.Auth
}

func NewMailService(cfg *config.Mail) *MailService {
	return &MailService{
		Config: cfg,
		Auth:   smtp.PlainAuth("", cfg.Mail, cfg.Pass, cfg.Addr),
	}
}

// from: %s\n Subject %s\n\n {body}
func (s *MailService) SendMail(ctx context.Context, to string, subject string, body []byte) error {
	return s.SendMails(ctx, []string{to}, subject, body)
}

func (s *MailService) SendMails(ctx context.Context, to []string, subject string, body []byte) error {
	log := s.log.With(slog.String("op", "smtp.MailService.SenvMail"))
	log.Debug("Start")

	stopCh := make(chan struct{}, 1)

	msg := fmt.Sprintf("From: %s\nSubject: %s\n\n%s", s.Config.From, subject, body)

	go func() {
		defer close(stopCh)

		err := smtp.SendMail(
			fmt.Sprintf("%s:%d", s.Config.Addr, s.Config.Port),
			s.Auth,
			s.Config.Mail,
			to,
			[]byte(msg),
		)
		if err != nil {
			log.Error("SendMail failed", slog.String("err", err.Error()))
			stopCh <- struct{}{}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stopCh:
		return nil
	}
}

// From %s\n
// To: %s\n
// Subject: %s\n\n 
// {body: string}\n\n 
// {filename: string, attachmentData: []byte}
func (s *MailService) SendMailWithAttachment(
	to string,
	subject string,
	body []byte,
	filename string,
	attachmentData []byte,
) error {
	msg := []byte(strings.Join([]string{
		"From: " + s.Config.From,
		"To: " + to,
		"Subject: " + subject,
		"MIME-version: 1.0",
		"Content-Type: multipart/mixed; boundary=\"boundary\"",
		"",
		"--boundary",
		"Content-Type: text/plain; charset=\"utf-8\"",
		"",
		string(body),
		"",
		"--boundary",
		"Content-Type: application/octet-stream; name=\"attachment\"",
		"Content-Transfer-Encoding: base64",
		"Content-Disposition: attachment; filename=\"" + filename + "\"",
		"",
		base64.StdEncoding.EncodeToString(attachmentData),
		"",
		"--boundary--",
	}, "\n"))

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.Config.Addr, s.Config.Port),
		s.Auth,
		s.Config.Mail,
		[]string{to},
		msg,
	)
}
func (s *MailService) SendMailsWithAttachment(
	to []string,
	subject string,
	body []byte,
	filename string,
	attachmentData []byte,
) error {
	buf := bytes.Buffer{}
	buf.WriteString("From: " + s.Config.From + "\n")
	buf.WriteString("Subject: " + subject + "\n\n")
	buf.WriteString("MIME-version: 1.0\n")
	buf.WriteString("Content-Type: multipart/mixed; boundary=\"boundary\"\n\n")

	buf.WriteString("--boundary\n")
	buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\n")
	buf.Write(body)
	buf.WriteString("\n\n")

	buf.WriteString("--boundary\n")
	buf.WriteString("Content-Type: application/octet-stream; name=\"attachment\"\n")
	buf.WriteString("Content-Transfer-Encoding: base64\n")
	buf.WriteString("Content-Disposition: attachment; filename=\"" + filename + "\"\n\n")
	buf.Write(attachmentData)
	buf.WriteString("\n\n")
	buf.WriteString("--boundary--")

	return s.SendMails(context.Background(), to, subject, buf.Bytes())
}
