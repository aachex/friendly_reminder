package reminder

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
	"github.com/artemwebber1/friendly_reminder/pkg/email"
)

// Reminder представляет собой объект, который в отдельной горутине
// присылает уведомления пользователям, подписанным на рассылку.
type Reminder interface {
	// StartSending в достаёт из базы данных электронные почты всех пользователей,
	// подписанных на рассылку, и отправляет им их списки дел c указанным интервалом.
	StartSending(ctx context.Context, d time.Duration)
}

type defaultReminder struct {
	sender    email.Sender // Для отправки электронных писем
	usersRepo repository.UsersRepository
	tasksRepo repository.TasksRepository
}

func New(s email.Sender, ur repository.UsersRepository, tr repository.TasksRepository) Reminder {
	return &defaultReminder{
		sender:    s,
		usersRepo: ur,
		tasksRepo: tr,
	}
}

// StartSending в достаёт из базы данных электронные почты всех пользователей,
// подписанных на рассылку, и отправляет им их списки дел c указанным интервалом.
func (s *defaultReminder) StartSending(ctx context.Context, d time.Duration) {
	for {
		log.Println("Sending emails")
		emails, err := s.usersRepo.GetEmailsSubscribed(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for _, email := range emails {
			go s.sendList(ctx, email)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(d):
			continue
		}
	}
}

func (s *defaultReminder) sendList(ctx context.Context, email string) {
	// Получаем список пользователя
	list, err := s.tasksRepo.GetList(ctx, email)
	if err != nil {
		log.Fatal(err)
	}

	// Преобразуем слайс list в строку вида:
	// 1. Задача 1
	// 2. Задача 2
	// ...
	body := "Here is your to-do list:"
	for i, item := range list {
		body += fmt.Sprintf("\n%d. %s", i+1, item.Value)
	}

	subject := "Your to-do list"
	if len(list) == 0 {
		// Отписываем пользователя от рассылки, если его список пуст.
		// Меняем заголовок и тело письма, чтобы уведомить пользователя об этом.
		subject = "You were unsubscribed from the mailing."
		body = "Your to-do list is empty. Add new tasks to your list and subscribe to the mailing."
		s.usersRepo.Subscribe(ctx, email, false) // Отписка от рассылки
	}

	if err = s.sender.Send(subject, body, email); err != nil {
		log.Fatal(err)
	}
}
