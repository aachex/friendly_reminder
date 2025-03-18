package main

import (
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
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}
	cfg := config.NewConfig(`config\config.json`)

	// Подключение к бд
	db, err := sql.Open(cfg.DbOptions.DriverName, cfg.DbOptions.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("PRAGMA FOREIGN_KEYS=ON")
	defer db.Close()

	// Инициализация логгера
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)

	// Инициализация репозиториев
	usersRepo := repository.NewUsersRepository(db)
	tasksRepo := repository.NewTasksRepository(db)
	unverifiedUsersRepo := repository.NewUnverifiedUsersRepository(db)

	// Объект для рассылки писем
	emailSender := email.NewSender(
		os.Getenv("EMAIL"),
		os.Getenv("EMAIL_PASSWORD"),
		cfg.EmailOptions.Host,
		cfg.EmailOptions.Port)

	// Создание контроллеров и добавление эндпоинтов
	mux := http.NewServeMux()
	usersController := controller.NewUsersController(usersRepo, unverifiedUsersRepo, emailSender, cfg)
	tasksController := controller.NewTasksController(tasksRepo)

	usersController.AddEndpoints(mux)
	tasksController.AddEndpoints(mux)

	// Запуск рассыльщика
	listSender := reminder.New(emailSender, usersRepo, tasksRepo)
	go listSender.StartSending(cfg.ListSenderOptions.Delay * time.Second)

	// Запуск сервера
	addr := ":" + cfg.Port
	fmt.Println("Listening:", cfg.Host+addr)

	serv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeoutInMilliseconds * time.Millisecond,
		WriteTimeout: cfg.WriteTimeoutInMilliseconds * time.Millisecond,
	}

	serv.ListenAndServe()
}
