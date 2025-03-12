package test

import (
	"slices"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

func TestAddUser(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	repo := repository.NewUsersRepository(db)

	_, err := repo.AddUser(mock.email, mock.pwd)

	if err != nil || !repo.EmailExists(mock.email) {
		t.Fatal(err)
	}
}

func TestMakeSigned(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	repo := repository.NewUsersRepository(db)

	_, err := repo.AddUser(mock.email, mock.pwd)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.MakeSigned(mock.email, true)
	if err != nil {
		t.Fatal(err)
	}

	signedEmails, err := repo.GetEmails()
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Contains(signedEmails, mock.email) {
		t.Fail()
	}
}
