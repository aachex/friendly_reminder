package email

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

type EmailSenderClient struct {
	from     string
	password string
	host     string
	port     string
	auth     smtp.Auth

	usersRepo repository.UsersRepository
	itemsRepo repository.ItemsRepository
}

func NewSender(from, password, host, port string, ur repository.UsersRepository, ir repository.ItemsRepository) *EmailSenderClient {
	return &EmailSenderClient{
		from:      from,
		password:  password,
		host:      host,
		port:      port,
		auth:      smtp.PlainAuth("", from, password, host),
		usersRepo: ur,
		itemsRepo: ir,
	}
}

func (s *EmailSenderClient) Send(subject, body, to string) error {
	msg := fmt.Appendf(nil, "Subject: %s\r\n%s", subject, body)

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
