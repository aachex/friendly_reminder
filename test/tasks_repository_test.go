package test

import (
	"context"
	"testing"

	repo "github.com/artemwebber1/friendly_reminder/internal/repository/sqlite"
)

func TestGetList(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	// Создаём пользователя
	usersRepo := repo.NewUsersRepository(db)
	const email = "abcde@gmail.com"
	const passwordHash = "hashedPassword"
	err := usersRepo.AddUser(t.Context(), email, passwordHash)

	if err != nil {
		t.Fatal(err)
	}

	// Пользователь добавляет новые дела в свой список
	itemsRepo := repo.NewTasksRepository(db)
	tasks := []string{"do homework", "smth", "##@@??"}

	for _, task := range tasks {
		itemsRepo.AddTask(t.Context(), task, email)
	}

	list, err := itemsRepo.GetList(t.Context(), email)
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

	tasksRepo := repo.NewTasksRepository(db)

	_, err = tasksRepo.AddTask(t.Context(), "error", "invalid@mail.com")
	if err == nil {
		t.Fatal(err)
	}
}

func TestAddTask_Timeout(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	tasksRepo := repo.NewTasksRepository(db)

	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()

	_, err := tasksRepo.AddTask(ctx, "error", "invalid@mail.com")
	if err != context.DeadlineExceeded {
		t.Fatalf("Wanted error %s, got %s", context.DeadlineExceeded, err)
	}
}
