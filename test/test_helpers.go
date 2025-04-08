package test

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/internal/controller"
	repo "github.com/artemwebber1/friendly_reminder/internal/repository/postgres"
	"github.com/artemwebber1/friendly_reminder/pkg/email"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver
)

// mock struct
type m struct {
	email, pwd string
}

var mock = m{
	email: "achex@mail.com",
	pwd:   "password4321",
}

var cfg *config.Config
var dbUsed config.DbConfig
var addr string

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	cfg = config.NewConfig("../config/config.json")
	dbUsed = cfg.Database.Postgres
	addr = cfg.Host + ":" + cfg.Port + cfg.Prefix
}

func openDb(t *testing.T) *sql.DB {
	db, err := sql.Open(dbUsed.DriverName, os.Getenv(dbUsed.ConnStrEnv))
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
	ur := repo.NewUsersRepository(db)
	uur := repo.NewUnverifiedUsersRepository(db)
	sender := getEmailSender(cfg.EmailOptions.Host, cfg.EmailOptions.Port)
	return controller.NewUsersController(ur, uur, sender, cfg)
}

func getTasksController(db *sql.DB) *controller.TasksController {
	tr := repo.NewTasksRepository(db)
	ur := repo.NewUsersRepository(db)
	return controller.NewTasksController(tr, ur, cfg)
}

func getEmailSender(emailHost, emailPort string) email.Sender {
	return email.NewSender(
		os.Getenv("EMAIL"),
		os.Getenv("EMAIL_PWD"),
		emailHost,
		emailPort)
}
