package test

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
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

var addr = cfg.Host + ":" + cfg.Port

func openDb(t *testing.T) *sql.DB {
	db, err := sql.Open(cfg.DbOptions.DriverName, cfg.DbOptions.DbPath)
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

func statusCodesMismatch(wanted, got int, body string) string {
	return fmt.Sprintf("Wanted status code %d, got %d.\nResponse body: %s", wanted, got, body)
}

func getJwt(t *testing.T, usersCtrl *controller.UsersController) string {
	resRec := httptest.NewRecorder()
	body := fmt.Appendf(nil, "{\"email\": \"%s\", \"password\": \"%s\"}", mock.email, mock.pwd)
	req, err := http.NewRequest(http.MethodPost, addr+"/login", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	usersCtrl.Login(resRec, req)
	if resRec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Wanted status code %d, got %d.\nResponse body: %s", http.StatusOK, resRec.Result().StatusCode, resRec.Body)
	}

	return resRec.Body.String()
}

func getUsersController(db *sql.DB) *controller.UsersController {
	ur := repository.NewUsersRepository(db)
	uur := repository.NewUnverifiedUsersRepository(db)
	sender := getEmailSender(cfg.EmailOptions.Host, cfg.EmailOptions.Port)
	return controller.NewUsersController(ur, uur, sender, cfg)
}

func getTasksController(db *sql.DB) *controller.TasksController {
	tr := repository.NewTasksRepository(db)
	return controller.NewTasksController(tr)
}

func getEmailSender(emailHost, emailPort string) email.Sender {
	return email.NewSender(
		os.Getenv("EMAIL"),
		os.Getenv("EMAIL_PASSWORD"),
		emailHost,
		emailPort)
}
