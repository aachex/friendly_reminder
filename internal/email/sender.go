package email

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

// Sender отправляет электронное письмо на указанный адрес.
type Sender interface {
	Send(subject, body, to string) error
}

type defaultSender struct {
	from     string
	password string
	host     string
	port     string
	auth     smtp.Auth

	usersRepo repository.UsersRepository
	tasksRepo repository.TasksRepository
}

func NewSender(from, password, host, port string, ur repository.UsersRepository, tr repository.TasksRepository) Sender {
	return &defaultSender{
		from:      from,
		password:  password,
		host:      host,
		port:      port,
		auth:      smtp.PlainAuth("", from, password, host),
		usersRepo: ur,
		tasksRepo: tr,
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

	log.Printf("Sent email from '%s' to '%s'\n", s.from, to)

	return err
}
