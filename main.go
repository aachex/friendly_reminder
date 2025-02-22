package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/internal/controller"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

const driverName = "sqlite3"
const dbPath = `D:\projects\golang\Web\friendly_reminder\db\database.db`

func main() {
	// Инициализировать конфиг
	config := config.NewConfig(`config\config.json`)

	// подключиться к бд
	db, err := sql.Open(driverName, dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Инициализировать контроллеры
	usersRepository := repository.NewUsersRepository(db)
	itemsRepository := repository.NewItemsRepository(db)

	id, _ := usersRepository.AddUser("achex@gmail.com", "password")
	itemsRepository.AddItem("value", id)

	usersController := controller.NewUsersController(usersRepository)

	// Добавить эндпоинты
	mux := http.NewServeMux()
	usersController.AddEndpoints(mux)

	// ...

	// Запустить сервер
	address := fmt.Sprintf("localhost:%d", config.Port)
	http.ListenAndServe(address, mux)
}
