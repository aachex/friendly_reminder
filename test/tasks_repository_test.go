package test

import (
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

func TestGetList(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	// Создаём пользователя
	repo := repository.NewUsersRepository(db)
	const email = "abcde@gmail.com"
	const passwordHash = "hashedPassword"
	_, err := repo.AddUser(email, passwordHash)

	if err != nil {
		t.Fatal(err)
	}

	// Пользователь добавляет новые дела в свой список
	itemsRepo := repository.NewTasksRepository(db)
	tasks := []string{"do homework", "smth", "##@@??"}

	for _, task := range tasks {
		itemsRepo.AddTask(task, email)
	}

	list, err := itemsRepo.GetList(email)
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 {
		t.Fatal("list is empty")
	}

	for i := range list {
		if list[i].Value != tasks[i] || list[i].UserEmail != email {
			t.Fatal("slices are not equal")
		}
	}
}

// Здесь тестируем обновление списка для несуществующего пользователя - должна возникнуть ошибка FOREIGN KEY constraint failed.
func TestAddTask_InvalidEmail(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	_, err := db.Exec("PRAGMA FOREIGN_KEYS=ON")
	if err != nil {
		t.Fatal(err)
	}

	tasksRepo := repository.NewTasksRepository(db)

	_, err = tasksRepo.AddTask("error", "invalid@mail.com")
	if err == nil {
		t.Fatal(err)
	}
}
