package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/internal/controller"
	"github.com/artemwebber1/friendly_reminder/internal/reminder"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
	"github.com/artemwebber1/friendly_reminder/pkg/email"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type App struct {
	cfg *config.Config
	srv *http.Server
}

func New(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) {
	// Подключение к бд
	db, err := sql.Open(a.cfg.DbOptions.DriverName, a.cfg.DbOptions.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("PRAGMA FOREIGN_KEYS=ON")
	defer db.Close()

	// Инициализация репозиториев
	usersRepo := repository.NewUsersRepository(db)
	tasksRepo := repository.NewTasksRepository(db)
	unverifiedUsersRepo := repository.NewUnverifiedUsersRepository(db)

	// Объект для рассылки писем
	emailSender := email.NewSender(
		os.Getenv("EMAIL"),
		os.Getenv("EMAIL_PWD"),
		a.cfg.EmailOptions.Host,
		a.cfg.EmailOptions.Port,
	)

	// Создание контроллеров и добавление эндпоинтов
	mux := http.NewServeMux()
	usersController := controller.NewUsersController(usersRepo, unverifiedUsersRepo, emailSender, a.cfg)
	tasksController := controller.NewTasksController(tasksRepo, a.cfg)

	usersController.AddEndpoints(mux)
	tasksController.AddEndpoints(mux)

	// Запуск рассыльщика
	listSender := reminder.New(emailSender, usersRepo, tasksRepo)
	go listSender.StartSending(ctx, a.cfg.ListSenderOptions.Delay*time.Second)

	// Запуск сервера
	addr := ":" + a.cfg.Port
	fmt.Println("Listening:", a.cfg.Host+addr)

	a.srv = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  a.cfg.ReadTimeoutInMilliseconds * time.Millisecond,
		WriteTimeout: a.cfg.WriteTimeoutInMilliseconds * time.Millisecond,
	}

	a.srv.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}
