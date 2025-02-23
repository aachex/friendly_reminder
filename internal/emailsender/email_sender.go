package emailsender

import (
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

type EmailSender struct {
	from     string
	password string
	host     string
	port     int

	usersRepo repository.UsersRepository
	itemsRepo repository.ItemsRepository
}

func New(from, password, host string, port int, ur repository.UsersRepository, ir repository.ItemsRepository) *EmailSender {
	return &EmailSender{
		from:      from,
		password:  password,
		host:      host,
		port:      port,
		usersRepo: ur,
		itemsRepo: ir,
	}
}

// StartMailing в отдельной горутине достаёт из базы данных электронные почты всех пользователей,
// подписанных на рассылку, и отправляет им их списки дел c указанным интервалом.
func (s *EmailSender) StartMailing(d time.Duration) {
	go func() {
		for {
			emails, err := s.usersRepo.GetEmails()
			if err != nil {
				log.Fatal(err)
			}

			for _, email := range emails {
				go func(email string) {
					// Получаем список пользователя
					list, err := s.itemsRepo.GetList(email)
					if err != nil {
						log.Fatal(err)
					}

					// Преобразуем слайс list в строку вида:
					// 1. Задача 1
					// 2. Задача 2
					// ...
					listStr := ""
					for _, item := range list {
						listStr += fmt.Sprintf("\n%d. %s", item.NumberInList, item.Value)
					}

					msg := []byte("Subject: A friendly reminder\nHere is your todo list" + listStr)
					if err = send(s.from, email, s.password, s.host, s.port, msg); err != nil {
						log.Fatal(err)
					}
				}(email)
			}

			time.Sleep(d)
		}
	}()
}

func send(from, to, password, host string, port int, msg []byte) error {
	auth := smtp.PlainAuth("", from, password, host)

	addr := fmt.Sprintf("%s:%d", host, port)
	err := smtp.SendMail(
		addr,
		auth,
		from,
		[]string{to},
		msg,
	)

	return err
}
