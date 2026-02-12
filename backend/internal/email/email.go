package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Service struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewService(host, port, username, password, from string) *Service {
	return &Service{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *Service) Send(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := strings.Join([]string{
		"From: " + s.from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		htmlBody,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
