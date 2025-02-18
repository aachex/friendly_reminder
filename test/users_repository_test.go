package test

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

func TestAddUser(t *testing.T) {
	db, err := sql.Open("sqlite3", `D:\projects\golang\Web\friendly_reminder\db\database.db`)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewUsersRepository(db)

	const email = "abcde@gmail.com"
	const passwordHash = "hashedPassword"
	repo.AddUser(email, passwordHash)

	emails, err := repo.GetEmails()
	if err != nil || !slices.Contains(emails, email) {
		t.Fatal(err)
	}

	if err = cleanDb(db); err != nil {
		t.Fatal(err)
	}
}

func cleanDb(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM Users;")
	return err
}
