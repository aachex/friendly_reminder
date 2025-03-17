package email

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
)

// Sender позволяет отправлять электронные письма с конкретного адреса.
type Sender interface {
	Send(subject, body, to string) error
}

type defaultSender struct {
	from     string
	password string
	host     string
	port     string
	auth     smtp.Auth
}

func NewSender(from, password, host, port string) Sender {
	return &defaultSender{
		from:     from,
		password: password,
		host:     host,
		port:     port,
		auth:     smtp.PlainAuth("", from, password, host),
	}
}

func (s *defaultSender) Send(subject, body, to string) error {
	msg := fmt.Appendf(
		nil,
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
			"Content-Transfer-Encoding: base64\r\n\r\n"+
			"%s",
		to, subject, base64.StdEncoding.EncodeToString([]byte(body)))

	addr := s.host + ":" + s.port
	err := smtp.SendMail(
		addr,
		s.auth,
		s.from,
		[]string{to},
		msg,
	)

	return err
}
