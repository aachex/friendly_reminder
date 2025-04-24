package test

import (
	"slices"
	"testing"

	repo "github.com/artemwebber1/friendly_reminder/internal/repository/postgres"
)

func TestAddUser(t *testing.T) {
	defer cleanDb(db, t)

	repo := repo.NewUsersRepository(db)

	err := repo.AddUser(t.Context(), mock.email, mock.pwd)

	if err != nil || !repo.EmailExists(t.Context(), mock.email) {
		t.Fatal(err)
	}
}

func TestMakeSigned(t *testing.T) {
	defer cleanDb(db, t)

	repo := repo.NewUsersRepository(db)

	err := repo.AddUser(t.Context(), mock.email, mock.pwd)
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

func TestGetByEmail(t *testing.T) {
	defer cleanDb(db, t)

	repo := repo.NewUsersRepository(db)
	err := repo.AddUser(t.Context(), mock.email, mock.pwd)
	if err != nil {
		t.Error(err)
	}

	u, err := repo.GetByEmail(t.Context(), mock.email)
	if err != nil {
		t.Error(err)
	}

	if u.Email != mock.email || u.Password != mock.pwd {
		t.Error("mismatch email or password")
	}
}
