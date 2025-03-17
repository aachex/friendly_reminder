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
	config := config.NewConfig(`config\config.json`)

	// Подключение к бд
	const driverName = "sqlite3"
	const dbPath = `D:\projects\golang\Web\friendly_reminder\db\database.db`
	db, err := sql.Open(driverName, dbPath)
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

	// Инициализация контроллеров и репозиториев
	usersRepo := repository.NewUsersRepository(db)
	tasksRepo := repository.NewTasksRepository(db)
	unverifiedUsersRepo := repository.NewUnverifiedUsersRepository(db)

	emailSender := email.NewSender(
		os.Getenv("EMAIL"),
		os.Getenv("EMAIL_PASSWORD"),
		config.EmailOptions.Host,
		config.EmailOptions.Port)

	// Создание контроллеров и добавление эндпоинтов
	mux := http.NewServeMux()
	usersController := controller.NewUsersController(usersRepo, unverifiedUsersRepo, emailSender, config)
	tasksController := controller.NewTasksController(tasksRepo)

	usersController.AddEndpoints(mux)
	tasksController.AddEndpoints(mux)

	// Запуск рассыльщика
	listSender := reminder.New(emailSender, usersRepo, tasksRepo)
	go listSender.StartSending(config.ListSenderOptions.DelayInSeconds * time.Second)

	// Запуск сервера
	addr := config.Host + ":" + config.Port
	fmt.Println("Listening:", addr)
	log.Fatal(http.ListenAndServe(":"+config.Port, mux))
}
