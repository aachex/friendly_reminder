package test

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

const email = "abcde@gmail.com"
const passwordHash = "hashedPassword"

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

	_, err = repo.AddUser(email, passwordHash)

	if err != nil || !repo.EmailExists(email) {
		t.Fatal(err)
	}
}

func TestMakeSigned(t *testing.T) {
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

	_, err = repo.AddUser(email, passwordHash)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.MakeSigned(email)
	if err != nil {
		t.Fatal(err)
	}

	signedEmails, err := repo.GetEmails()
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Contains(signedEmails, email) {
		t.Fail()
	}
}
