package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/artemwebber1/friendly_reminder/internal/email"
	"github.com/artemwebber1/friendly_reminder/internal/models"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

// Возможные ошибки

const (
	userAlreadyExists = "Пользователь с данной электронной почтой уже существует"
)

type UsersController struct {
	repo        repository.UsersRepository
	emailSender email.EmailSenderClient
}

func NewUsersController(repo repository.UsersRepository) *UsersController {
	return &UsersController{
		repo: repo,
	}
}

func (c *UsersController) AddEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /new-user", c.AddUser)
	mux.HandleFunc("POST /user-auth", c.AuthUser)
}

// AddUser создаёт нового пользователя в базе данных.
//
// Обрабатывает POST запросы по пути '/new-user'.
func (c *UsersController) AddUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var user models.User
	json.Unmarshal(bodyBytes, &user)

	if c.repo.EmailExists(user.Email) {
		http.Error(w, userAlreadyExists, http.StatusForbidden)
		return
	}

	// Все проверки прошли успешно, отправляем пользователю на почту ссылку для подтверждения электронной почты

	confirmLink := "..."

	const subject = "Подтверждение электронной почты"
	body := fmt.Sprintf("Пожалуйста, подтвердите электронную почту, перейдя по ссылке:\n%s\n\nЕсли вы не ждали этого письма, проигнорируйте его.", confirmLink)
	c.emailSender.Send(
		subject,
		body,
		user.Email)
}

// SignUser подписывает пользователя с указанным email на рассылку писем.
//
// Обрабатывает PATCH запросы по пути '/sign-user?{email}'.
func (c *UsersController) SignUser(w http.ResponseWriter, r *http.Request) {

}

// AuthUser осуществляет вход уже существующего пользователя в систему.
//
// Обрабатывает POST запросы по пути '/user-auth'.
func (C *UsersController) AuthUser(w http.ResponseWriter, r *http.Request) {
	// Получить эл. почту и пароль

	// Проверить, что они корректны. Если не корректны, вернуть код 403.

	// Создать jwt и вернуть его
}
