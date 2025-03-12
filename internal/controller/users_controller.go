package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/internal/hasher"
	"github.com/artemwebber1/friendly_reminder/internal/models"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
	"github.com/artemwebber1/friendly_reminder/pkg/email"
	"github.com/artemwebber1/friendly_reminder/pkg/middleware"
)

type UsersController struct {
	emailSender email.Sender
	cfg         *config.Config

	usersRepo           repository.UsersRepository
	unverifiedUsersRepo repository.UnverifiedUsersRepository
}

func NewUsersController(
	ur repository.UsersRepository,
	uur repository.UnverifiedUsersRepository,
	emailSender email.Sender,
	cfg *config.Config) *UsersController {
	return &UsersController{
		usersRepo:           ur,
		unverifiedUsersRepo: uur,
		emailSender:         emailSender,
		cfg:                 cfg,
	}
}

func (c *UsersController) AddEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /new-user", c.SendConfirmEmailLink)
	mux.HandleFunc("POST /user-auth", c.AuthUser)
	mux.HandleFunc("GET /confirm-email", c.ConfirmEmail)
	mux.HandleFunc("PATCH /sign-user", middleware.RequireAuthorization(c.SignUser))
}

// AddUser создаёт нового пользователя в базе данных.
//
// Обрабатывает POST запросы по пути '/new-user'.
func (c *UsersController) SendConfirmEmailLink(w http.ResponseWriter, r *http.Request) {
	user, err := readBody[models.User](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if c.usersRepo.EmailExists(user.Email) {
		http.Error(w, "Пользователь с данной электронной почтой уже существует", http.StatusForbidden)
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
	confirmLink := c.cfg.Host + ":" + c.cfg.Port + "/confirm-email?t=" + confirmToken

	log.Printf("Sending confirm email link to '%s'...\n", user.Email)

	const subject = "Подтверждение электронной почты"
	body := fmt.Sprintf("Пожалуйста, подтвердите свою электронную почту, перейдя по ссылке:\n%s\n\nЕсли вы не запрашивали это письмо, проигнорируйте его.", confirmLink)
	c.emailSender.Send(
		subject,
		body,
		user.Email)
	w.Write([]byte(confirmToken))
}

// ConfirmEmail является эндпоинтом, на который пользователь попадёт, подтверждая электронную почту.
//
// Обрабатывает GET запросы по пути '/confirm-email'.
func (c *UsersController) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("t")
	if !c.unverifiedUsersRepo.TokenExists(token) {
		http.Error(w, "Недействительный токен", http.StatusForbidden)
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
// Обрабатывает PATCH запросы по пути '/sign-user'.
func (c *UsersController) SignUser(w http.ResponseWriter, r *http.Request) {
	usrEmail := r.URL.Query().Get("email") // TODO: достать email из jwt токена
	sign, err := strconv.ParseBool(r.URL.Query().Get("sign"))
	if err != nil {
		http.Error(w, "Параметр sign: неправильное значение", http.StatusBadRequest)
		return
	}
	c.usersRepo.MakeSigned(usrEmail, sign)
}

// AuthUser осуществляет вход уже существующего пользователя в систему.
//
// Обрабатывает POST запросы по пути '/user-auth'.
func (C *UsersController) AuthUser(w http.ResponseWriter, r *http.Request) {
	// email := r.Body.Get("email")
	// pwd := r.Body.Get("password")
	//
	// if emailExists(email) -> http.Error(403)
	//
	// jwtHeaders := { alg: "sha256" }
	// jwtPayload := { email }
	// jwt := NewJwt(jwtHeaders, jwtPayload)
	//
	// w.Write(jwt)
}
