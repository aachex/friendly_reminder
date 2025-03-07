package reminder

import (
	"fmt"
	"log"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/email"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

type ListSender struct {
	sender    email.EmailSenderClient
	usersRepo repository.UsersRepository
	itemsRepo repository.ItemsRepository
}

func New(s email.EmailSenderClient, ur repository.UsersRepository, ir repository.ItemsRepository) *ListSender {
	return &ListSender{
		sender:    s,
		usersRepo: ur,
		itemsRepo: ir,
	}
}

// StartMailing в отдельной горутине достаёт из базы данных электронные почты всех пользователей,
// подписанных на рассылку, и отправляет им их списки дел c указанным интервалом.
func (s *ListSender) StartSending(d time.Duration) {
	for {
		log.Println("Sending emails")
		emails, err := s.usersRepo.GetEmails()
		if err != nil {
			log.Fatal(err)
		}

		for _, email := range emails {
			go s.sendList(email)
		}

		time.Sleep(d)
	}
}

func (s *ListSender) sendList(email string) {
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

	subject := "Ваш список дел"
	if len(list) == 0 {
		// Отписываем пользователя от рассылки, если его список пуст, и информируем его об этом.
		subject = "Вы были отписаны от рассылки"
		listStr = "Ваш список дел пуст. Вы будете отписаны от рассылки, пока не добавите новые дела и не подпишетесь на рассылку снова."
		s.usersRepo.MakeSigned(email, false) // Отписка от рассылки
	}

	if err = s.sender.Send(subject, listStr, email); err != nil {
		log.Fatal(err)
	}
}
