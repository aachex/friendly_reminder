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
	"github.com/artemwebber1/friendly_reminder/internal/email"
	"github.com/artemwebber1/friendly_reminder/internal/reminder"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
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
	defer db.Close()

	// Инициализация логгера
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)

	// Инициализация контроллеров и репозиториев
	usersRepository := repository.NewUsersRepository(db)
	itemsRepository := repository.NewItemsRepository(db)

	usersController := controller.NewUsersController(usersRepository)

	// Добавление эндпоинтов
	mux := http.NewServeMux()
	usersController.AddEndpoints(mux)

	usersRepository.AddUser("chekhonin.artem@gmail.com", "aaaa")
	usersRepository.MakeSigned("chekhonin.artem@gmail.com", true)
	itemsRepository.AddItem("погладить кошку", "chekhonin.artem@gmail.com")
	itemsRepository.AddItem("сделать дз", "chekhonin.artem@gmail.com")

	// Запуск рассыльщика
	emailSender := email.NewSender(
		os.Getenv("EMAIL"),
		os.Getenv("EMAIL_PASSWORD"),
		config.EmailHost,
		config.EmailPort,
		usersRepository,
		itemsRepository)

	listSender := reminder.New(*emailSender, usersRepository, itemsRepository)
	go listSender.StartSending(60 * time.Second)

	// Запуск сервера
	fmt.Println("Listening:", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, mux))
}
