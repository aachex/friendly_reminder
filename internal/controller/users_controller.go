package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/internal/hasher"
	"github.com/artemwebber1/friendly_reminder/internal/models"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
	"github.com/artemwebber1/friendly_reminder/pkg/email"
	"github.com/artemwebber1/friendly_reminder/pkg/middleware"
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
	mux.HandleFunc("POST /new-user", c.SendConfirmEmailLink)
	mux.HandleFunc("POST /login", c.Login)
	mux.HandleFunc("GET /confirm-email", c.ConfirmEmail)
	mux.HandleFunc("PATCH /sign-user", middleware.RequireAuthorization(c.SignUser))
}

// AddUser создаёт нового пользователя в базе данных.
//
// Обрабатывает POST запросы по пути '/new-user'.
func (c *UsersController) SendConfirmEmailLink(w http.ResponseWriter, r *http.Request) {
	user, err := readBody[models.User](r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
	}

	if c.usersRepo.EmailExists(user.Email) {
		http.Error(w, "User with this email already exists", http.StatusForbidden)
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
		http.Error(w, "Error creating confirm token", http.StatusInternalServerError)
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
		http.Error(w, "Invalid confirm token", http.StatusForbidden)
		return
	}

	user, err := c.unverifiedUsersRepo.GetUserByToken(token)
	if err != nil {
		http.Error(w, "Impossible to confirm email: undefined user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Почта подтверждена"))
	err = c.unverifiedUsersRepo.DeleteToken(token)
	if err != nil {
		http.Error(w, "Failed to delete confirm token", http.StatusForbidden)
		return
	}

	// Пользователь успешно подтвердил электронную почту, добавляем его в базу данных
	c.usersRepo.AddUser(user.Email, user.Password)
	w.WriteHeader(http.StatusCreated)
}

var key = []byte(os.Getenv("SECRET_STR"))

// SignUser подписывает пользователя с указанным email на рассылку писем.
//
// Обрабатывает PATCH запросы по пути '/sign-user'.
func (c *UsersController) SignUser(w http.ResponseWriter, r *http.Request) {
	rawTok := r.Header.Get("Authorization")
	if len(rawTok) > 7 {
		rawTok = rawTok[7:] // Отрезаем часть "Bearer: "
	}

	jwtClaims, err := readJWT(rawTok, key)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}

	userEmail, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, "Invalid field 'sub' in token claims", http.StatusForbidden)
		return
	}

	sign, err := strconv.ParseBool(r.URL.Query().Get("sign"))
	if err != nil {
		http.Error(w, "Invalid value for 'sign' param", http.StatusBadRequest)
		return
	}
	c.usersRepo.MakeSigned(userEmail, sign)
}

// Login осуществляет вход уже существующего пользователя в систему.
//
// Обрабатывает POST запросы по пути '/login'.
func (c *UsersController) Login(w http.ResponseWriter, r *http.Request) {
	user, err := readBody[models.User](r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if !c.usersRepo.UserExists(user.Email, hasher.Hash(user.Password)) {
		http.Error(w, "Invalid email or password", http.StatusForbidden)
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
		http.Error(w, "Error signing token", http.StatusBadRequest)
		return
	}

	w.Write([]byte(tokStr))
}
