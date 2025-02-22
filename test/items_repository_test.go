package test

import (
	"database/sql"
	"testing"

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
	userId, err := repo.AddUser(email, passwordHash)

	if err != nil {
		t.Fatal(err)
	}

	// Пользователь добавляет новые дела в свой список
	itemsRepo := repository.NewItemsRepository(db)
	const task = "Сделать дз"
	itemsRepo.AddItem(task, userId)

	list, err := itemsRepo.GetList(userId)
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 || list[0].Value != task || list[0].NumberInList != int64(len(list)) {
		t.Fail()
	}
}
