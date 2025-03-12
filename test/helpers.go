package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/internal/controller"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
	"github.com/artemwebber1/friendly_reminder/pkg/email"
)

// mock struct
type m struct {
	email, pwd string
}

var mock = m{
	email: "achex@mail.com",
	pwd:   "password4321",
}

var cfg = config.NewConfig(`D:\projects\golang\Web\friendly_reminder\config\config.json`)

func openDb(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", `D:\projects\golang\Web\friendly_reminder\db\database.db`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func cleanDb(db *sql.DB, t *testing.T) {
	_, err := db.Exec("DELETE FROM tasks; DELETE FROM users; DELETE FROM unverified_users;")
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
}

func getUsersController(db *sql.DB) controller.UsersController {
	ur := repository.NewUsersRepository(db)
	uur := repository.NewUnverifiedUsersRepository(db)
	sender := getEmailSender(cfg.EmailOptions.Host, cfg.EmailOptions.Port)
	return *controller.NewUsersController(ur, uur, sender, cfg)
}

func getEmailSender(emailHost, emailPort string) email.Sender {
	return email.NewSender(
		os.Getenv("EMAIL"),
		os.Getenv("EMAIL_PASSWORD"),
		emailHost,
		emailPort)
}
