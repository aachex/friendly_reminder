package test

import (
	"database/sql"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

func TestAddUser(t *testing.T) {
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

	repo := repository.NewUsersRepository(db)

	const email = "abcde@gmail.com"
	const passwordHash = "hashedPassword"
	_, err = repo.AddUser(email, passwordHash)

	if err != nil || !repo.EmailExists(email) {
		t.Fatal(err)
	}
}
