package test

import (
	"database/sql"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/models"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

func TestGetList(t *testing.T) {
	db, err := sql.Open("sqlite3", `D:\projects\golang\Web\friendly_reminder\db\database.db`)
	if err != nil {
		t.Fatal(err)
	}
	defer func(db *sql.DB) {
		if err = cleanDb(db); err != nil {
			t.Fatal(err)
		}
		db.Close()
	}(db)

	// Создаём пользователя
	repo := repository.NewUsersRepository(db)
	const email = "abcde@gmail.com"
	const passwordHash = "hashedPassword"
	_, err = repo.AddUser(email, passwordHash)

	if err != nil {
		t.Fatal(err)
	}

	// Пользователь добавляет новые дела в свой список
	itemsRepo := repository.NewTasksRepository(db)
	tasks := []models.Task{
		{Value: "сделать дз", NumberInList: 1},
		{Value: "погладить кошку", NumberInList: 2},
		{Value: "исправить оценки", NumberInList: 3},
	}

	for _, task := range tasks {
		itemsRepo.AddTask(task.Value, email)
	}

	list, err := itemsRepo.GetList(email)
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 {
		t.Fatal("list is empty")
	}

	for i := range list {
		if list[i].Value != tasks[i].Value || list[i].NumberInList != tasks[i].NumberInList || list[i].UserEmail != email {
			t.Fatal("slices are not equal")
		}
	}

	if list[len(list)-1].NumberInList != int64(len(list)) {
		t.Fatal("invalid numeration")
	}
}

// Здесь тестируем обновление списка для несуществующего пользователя - должна возникнуть ошибка FOREIGN KEY constraint failed.
func TestAddTask_InvalidEmail(t *testing.T) {
	db, err := sql.Open("sqlite3", `D:\projects\golang\Web\friendly_reminder\db\database.db`)
	if err != nil {
		t.Fatal(err)
	}
	defer func(db *sql.DB) {
		if err = cleanDb(db); err != nil {
			t.Fatal(err)
		}
		db.Close()
	}(db)

	_, err = db.Exec("PRAGMA FOREIGN_KEYS=ON")
	if err != nil {
		t.Fatal(err)
	}

	tasksRepo := repository.NewTasksRepository(db)

	err = tasksRepo.AddTask("error", "invalid@mail.com")
	if err == nil {
		t.Fail()
	}
}
