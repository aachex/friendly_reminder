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
	itemsRepo := repository.NewItemsRepository(db)
	tasks := []models.ListItem{
		{Value: "сделать дз", NumberInList: 1},
		{Value: "погладить кошку", NumberInList: 2},
		{Value: "исправить оценки", NumberInList: 3},
	}

	for _, task := range tasks {
		itemsRepo.AddItem(task.Value, email)
	}

	list, err := itemsRepo.GetList(email)
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 || !equalSlices(list, tasks) || list[len(list)-1].NumberInList != int64(len(list)) {
		t.Fail()
	}
}
