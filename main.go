package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/internal/controller"
	"github.com/artemwebber1/friendly_reminder/internal/emailsender"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

func main() {
	config := config.NewConfig(`config\config.json`)

	// Подключение к бд
	const driverName = "sqlite3"
	const dbPath = `D:\projects\golang\Web\friendly_reminder\db\database.db`
	db, err := sql.Open(driverName, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация контроллеров и репозиториев
	usersRepository := repository.NewUsersRepository(db)
	itemsRepository := repository.NewItemsRepository(db)

	// usersRepository.AddUser(config.Email, "aaaa")
	// itemsRepository.AddItem("Сделать дз", config.Email)
	// itemsRepository.AddItem("Погладить кота", config.Email)
	// itemsRepository.AddItem("Ботать олимпиады", config.Email)

	usersController := controller.NewUsersController(usersRepository)

	// Добавление эндпоинтов
	mux := http.NewServeMux()
	usersController.AddEndpoints(mux)

	// Запуск рассыльщика
	emailSender := emailsender.New(config.Email, config.EmailPassword, config.EmailHost, config.EmailPort, usersRepository, itemsRepository)
	emailSender.StartMailing(15 * time.Second)

	// Запуск сервера
	address := fmt.Sprintf(":%d", config.Port)
	fmt.Println("Listening:", address)
	log.Fatal(http.ListenAndServe(address, mux))
}
