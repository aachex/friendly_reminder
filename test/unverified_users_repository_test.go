package test

import (
	"database/sql"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

func TestCreateToken(t *testing.T) {
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

	tokRepo := repository.NewUnverifiedUsersRepository(db)
	tok, err := tokRepo.CreateToken(email, passwordHash)
	if err != nil {
		t.Fatal(err)
	}

	if !tokRepo.TokenExists(tok) || !tokRepo.HasToken(email) {
		t.Fail()
	}
}

func TestGetUserByToken(t *testing.T) {
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

	tokRepo := repository.NewUnverifiedUsersRepository(db)

	tok, err := tokRepo.CreateToken(email, passwordHash)
	if err != nil {
		t.Fatal(err)
	}

	if !tokRepo.HasToken(email) {
		t.Fatalf("У пользователя %s нет токена", email)
	}

	user, err := tokRepo.GetUserByToken(tok)
	if err != nil {
		t.Fatal(err)
	}

	if user.Email != email || user.Password != passwordHash {
		t.Fatal("Несоответствие эл. почты или пароля")
	}
}
