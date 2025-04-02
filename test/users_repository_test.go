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

	_, err := repo.AddUser(t.Context(), mock.email, mock.pwd)

	if err != nil || !repo.EmailExists(t.Context(), mock.email) {
		t.Fatal(err)
	}
}

func TestMakeSigned(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	repo := repository.NewUsersRepository(db)

	_, err := repo.AddUser(t.Context(), mock.email, mock.pwd)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Subscribe(t.Context(), mock.email, true)
	if err != nil {
		t.Fatal(err)
	}

	signedEmails, err := repo.GetEmailsSubscribed(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Contains(signedEmails, mock.email) {
		t.Fail()
	}
}
