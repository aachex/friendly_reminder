package test

import (
	"testing"

	repo "github.com/artemwebber1/friendly_reminder/internal/repository/sqlite"
)

func TestCreateToken(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	tokRepo := repo.NewUnverifiedUsersRepository(db)
	tok, err := tokRepo.CreateToken(mock.email, mock.pwd)
	if err != nil {
		t.Fatal(err)
	}

	if !tokRepo.TokenExists(tok) || !tokRepo.HasToken(mock.email) {
		t.Fail()
	}
}

func TestGetUserByToken(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	tokRepo := repo.NewUnverifiedUsersRepository(db)

	tok, err := tokRepo.CreateToken(mock.email, mock.pwd)
	if err != nil {
		t.Fatal(err)
	}

	if !tokRepo.HasToken(mock.email) {
		t.Fatalf("Token doesn't exist: %s", mock.email)
	}

	user, err := tokRepo.GetUserByToken(tok)
	if err != nil {
		t.Fatal(err)
	}

	if user.Email != mock.email || user.Password != mock.pwd {
		t.Fatal("Email and password mismatch")
	}
}
