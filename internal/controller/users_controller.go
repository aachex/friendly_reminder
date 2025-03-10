package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/internal/email"
	"github.com/artemwebber1/friendly_reminder/internal/hasher"
	"github.com/artemwebber1/friendly_reminder/internal/models"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

const (
	userAlreadyExists = "Пользователь с данной электронной почтой уже существует"
	invalidToken      = "Недействительный токен"
) // Возможные ошибки

type UsersController struct {
	usersRepo           repository.UsersRepository
	unverifiedUsersRepo repository.UnverifiedUsersRepository
	emailSender         email.Sender
	config              config.Config
}

func NewUsersController(
	ur repository.UsersRepository,
	tr repository.UnverifiedUsersRepository,
	emailSender email.Sender,
	cfg config.Config) *UsersController {
	return &UsersController{
		usersRepo:           ur,
		unverifiedUsersRepo: tr,
		emailSender:         emailSender,
		config:              cfg,
	}
}

func (c *UsersController) AddEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /new-user", c.AddUser)
	mux.HandleFunc("POST /user-auth", c.AuthUser)
	mux.HandleFunc("GET /confirm-email", c.ConfirmEmail)
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

	if c.usersRepo.EmailExists(user.Email) {
		http.Error(w, userAlreadyExists, http.StatusForbidden)
		return
	}

	// Отправляем пользователю на почту ссылку для подтверждения электронной почты
	var confirmToken string
	if !c.unverifiedUsersRepo.HasToken(user.Email) {
		confirmToken, err = c.unverifiedUsersRepo.CreateToken(user.Email, hasher.Hash([]byte(user.Password)))
	} else {
		confirmToken, err = c.unverifiedUsersRepo.UpdateToken(user.Email)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ссылка для подтверждения электронной почты
	confirmLink := c.config.Host + ":" + c.config.Port + "/confirm-email?token=" + confirmToken

	const subject = "Подтверждение электронной почты"
	body := fmt.Sprintf("Пожалуйста, подтвердите свою электронную почту, перейдя по ссылке:\n%s\n\nЕсли вы не запрашивали это письмо, проигнорируйте его.", confirmLink)
	c.emailSender.Send(
		subject,
		body,
		user.Email)
}

// ConfirmEmail является эндпоинтом, на который пользователь попадёт, подтверждая электронную почту.
//
// Обрабатывает GET запросы по пути '/confirm-email?{token}'.
func (c *UsersController) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if !c.unverifiedUsersRepo.TokenExists(token) {
		http.Error(w, invalidToken, http.StatusForbidden)
		return
	}

	user, err := c.unverifiedUsersRepo.GetUserByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Почта подтверждена"))
	err = c.unverifiedUsersRepo.DeleteToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Пользователь успешно подтвердил электронную почту, добавляем его в базу данных
	c.usersRepo.AddUser(user.Email, user.Password)
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
