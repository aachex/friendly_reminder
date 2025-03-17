package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/internal/hasher"
	"github.com/artemwebber1/friendly_reminder/internal/models"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
	"github.com/artemwebber1/friendly_reminder/pkg/email"
	mw "github.com/artemwebber1/friendly_reminder/pkg/middleware"
	"github.com/golang-jwt/jwt/v5"
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
	mux.HandleFunc(
		"POST /new-user",
		c.SendConfirmEmailLink)

	mux.HandleFunc(
		"POST /login",
		c.Login)

	mux.HandleFunc(
		"GET /confirm-email",
		c.ConfirmEmail)

	mux.HandleFunc(
		"PATCH /subscribe",
		mw.UseAuthorization(c.SubscribeUser))
}

// AddUser создаёт нового пользователя в базе данных.
//
// Обрабатывает POST запросы по пути '/new-user'.
func (c *UsersController) SendConfirmEmailLink(w http.ResponseWriter, r *http.Request) {
	user, err := readBody[models.User](r.Body)
	if err != nil {
		http.Error(w, errReadingBody.Error(), http.StatusBadRequest)
	}

	if c.usersRepo.EmailExists(user.Email) {
		http.Error(w, "user with this email already exists", http.StatusForbidden)
		return
	}

	// Отправляем пользователю на почту ссылку для подтверждения электронной почты

	var confirmToken string
	if !c.unverifiedUsersRepo.HasToken(user.Email) {
		confirmToken, err = c.unverifiedUsersRepo.CreateToken(user.Email, hasher.Hash(user.Password))
	} else {
		confirmToken, err = c.unverifiedUsersRepo.UpdateToken(user.Email)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("error creating confirm token: %s", err), http.StatusInternalServerError)
		return
	}

	// Ссылка для подтверждения электронной почты
	confirmLink := c.cfg.Host + ":" + c.cfg.Port + "/confirm-email?t=" + confirmToken

	log.Printf("Sending an email confirmation link to '%s'...\n", user.Email)

	const subject = "Email confirmation"
	body := fmt.Sprintf("Please, confirm your email by clicking on the link:\n%s\n\nIf you didn't request this mail, ignore it.", confirmLink)
	c.emailSender.Send(
		subject,
		body,
		user.Email)
	w.Write([]byte(confirmToken))

	log.Printf("Sent an email confirmation link to '%s'\n", user.Email)
}

// ConfirmEmail является эндпоинтом, на который пользователь попадёт, подтверждая электронную почту.
//
// Обрабатывает GET запросы по пути '/confirm-email'.
func (c *UsersController) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("t")
	if !c.unverifiedUsersRepo.TokenExists(token) {
		http.Error(w, "invalid confirm token", http.StatusForbidden)
		return
	}

	user, err := c.unverifiedUsersRepo.GetUserByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Email confirmed succesfully"))
	err = c.unverifiedUsersRepo.DeleteToken(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete confirm token: %s", err), http.StatusForbidden)
		return
	}

	// Пользователь успешно подтвердил электронную почту, добавляем его в базу данных
	c.usersRepo.AddUser(user.Email, user.Password)
	w.WriteHeader(http.StatusCreated)
}

// SubscribeUser подписывает пользователя с указанным email на рассылку писем.
//
// Обрабатывает PATCH запросы по пути '/subscribe'.
func (c *UsersController) SubscribeUser(w http.ResponseWriter, r *http.Request) {
	rawJwt := getRawJwtFromHeader(r.Header)
	jwtClaims, err := readJWT(rawJwt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	userEmail, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	subscribe, err := strconv.ParseBool(r.URL.Query().Get("subscribe"))
	if err != nil {
		http.Error(w, "invalid value for 'subscribe' param", http.StatusBadRequest)
		return
	}
	c.usersRepo.Subscribe(userEmail, subscribe)
}

// Login осуществляет вход уже существующего пользователя в систему.
//
// Обрабатывает POST запросы по пути '/login'.
func (c *UsersController) Login(w http.ResponseWriter, r *http.Request) {
	user, err := readBody[models.User](r.Body)
	if err != nil {
		http.Error(w, errReadingBody.Error(), http.StatusBadRequest)
		return
	}

	if !c.usersRepo.UserExists(user.Email, hasher.Hash(user.Password)) {
		http.Error(w, "invalid email or password", http.StatusForbidden)
		return
	}

	// Создание jwt
	claims := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokStr, err := tok.SignedString(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(tokStr))
}
