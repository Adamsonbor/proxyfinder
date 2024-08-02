package smtpmail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"proxyfinder/internal/config"
	"strings"
)

type MailProvider interface {
	SendMail(to string, subject string, body string) error
	SendMails(to []string, subject string, body string) error
	SendMailWithAttachment(to string, subject string, body string, filename string, attachmentData []byte) error
	SendMailsWithAttachment(to []string, subject string, body string, filename string, attachmentData []byte) error
}

type SMTPProvider struct {
	Config *config.Mail
	Auth   smtp.Auth
}

func NewSMTPProvider(cfg *config.Mail) *SMTPProvider {
	return &SMTPProvider{
		Config: cfg,
		Auth:   smtp.PlainAuth("", cfg.Mail, cfg.Pass, cfg.Addr),
	}
}

func (s *SMTPProvider) SendMail(to string, subject string, body []byte) error {
	return s.SendMails([]string{to}, subject, body)
}

func (s *SMTPProvider) SendMails(to []string, subject string, body []byte) error {
	msg := fmt.Sprintf("From: %s\nSubject: %s\n\n%s", s.Config.From, subject, body)

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.Config.Addr, s.Config.Port),
		s.Auth,
		s.Config.Mail,
		to,
		[]byte(msg),
	)
}

func (s *SMTPProvider) SendMailWithAttachment(
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
func (s *SMTPProvider) SendMailsWithAttachment(
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

	return s.SendMails(to, subject, buf.Bytes())
}
