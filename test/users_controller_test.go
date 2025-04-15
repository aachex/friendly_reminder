package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/hasher"
	repo "github.com/artemwebber1/friendly_reminder/internal/repository/postgres"
)

func TestSendConfirmEmailLink(t *testing.T) {
	defer cleanDb(db, t)

	usersCtrl := getUsersController(db)

	body := bytes.NewReader(fmt.Appendf(nil, "{\"email\": \"%s\", \"password\": \"%s\"}", mock.email, mock.pwd))
	req, err := http.NewRequest(http.MethodPost, addr+"/new-user", body)
	if err != nil {
		t.Fatal(err)
	}

	resRec := httptest.NewRecorder()
	usersCtrl.SendConfirmEmailLink(resRec, req)

	if resRec.Result().StatusCode != http.StatusOK {
		t.Fatal(statusCodesMismatch(http.StatusOK, resRec.Result().StatusCode, resRec.Body.String()))
	}

	tok := resRec.Body.String()

	uur := repo.NewUnverifiedUsersRepository(db)
	if !uur.HasToken(mock.email) || !uur.TokenExists(tok) {
		t.Fatalf("Invalid token: %s", tok)
	}
}

func TestConfirmEmail(t *testing.T) {
	defer cleanDb(db, t)

	body := bytes.NewReader(fmt.Appendf(nil, "{\"email\": \"%s\", \"password\": \"%s\"}", mock.email, mock.pwd))
	req, err := http.NewRequest("POST", addr+"/new-user", body)
	if err != nil {
		t.Fatal(err)
	}
	resRec := httptest.NewRecorder()

	usersCtrl := getUsersController(db)
	usersCtrl.SendConfirmEmailLink(resRec, req)

	tok := resRec.Body.String()

	uur := repo.NewUnverifiedUsersRepository(db)
	if !uur.HasToken(mock.email) || !uur.TokenExists(tok) {
		t.Fatalf("Invalid token: %s", tok)
	}

	req, err = http.NewRequest(http.MethodGet, addr+"/confirm-email?t="+tok, nil)
	if err != nil {
		t.Fatal(err)
	}
	resRec = httptest.NewRecorder()
	usersCtrl.ConfirmEmail(resRec, req)

	if uur.TokenExists(tok) {
		t.Fatal("Failed to create new user")
	}
}

func TestSubscribeUser_Unauthorized(t *testing.T) {
	defer cleanDb(db, t)

	resRec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPatch, addr+"/subscribe?subscribe=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	usersCtrl := getUsersController(db)
	usersCtrl.SubscribeUser(resRec, req)
	if resRec.Result().StatusCode != http.StatusForbidden {
		t.Fatal(statusCodesMismatch(http.StatusForbidden, resRec.Result().StatusCode, resRec.Body.String()))
	}
}

func TestSubscribeUser(t *testing.T) {
	defer cleanDb(db, t)

	usersCtrl := getUsersController(db)

	// Сначала регистрируем пользователя для получения токена авторизации
	usersRepo := repo.NewUsersRepository(db)
	usersRepo.AddUser(t.Context(), mock.email, hasher.Hash(mock.pwd))

	tok := getJwt(t, usersCtrl)

	req, err := http.NewRequest(http.MethodPatch, addr+"/users/subscribe?subscribe=true", nil)
	req.Header.Add("Authorization", "Bearer "+tok)
	if err != nil {
		t.Fatal(err)
	}

	resRec := httptest.NewRecorder()
	usersCtrl.SubscribeUser(resRec, req)
	if resRec.Result().StatusCode != http.StatusOK {
		t.Fatal(statusCodesMismatch(http.StatusOK, resRec.Result().StatusCode, resRec.Body.String()))
	}
}

func TestLogin(t *testing.T) {
	defer cleanDb(db, t)

	usersRepo := repo.NewUsersRepository(db)
	err := usersRepo.AddUser(t.Context(), mock.email, hasher.Hash(mock.pwd))
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	resRec := httptest.NewRecorder()
	body := fmt.Appendf(nil, "{\"email\": \"%s\", \"password\": \"%s\"}", mock.email, mock.pwd)
	req, err := http.NewRequest(http.MethodPost, addr+"/users/login", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	usersCtrl := getUsersController(db)
	usersCtrl.Login(resRec, req)
	if resRec.Result().StatusCode != http.StatusOK {
		t.Fatal(statusCodesMismatch(http.StatusOK, resRec.Result().StatusCode, resRec.Body.String()))
	}
}

func TestLogin_UserDoesntExist(t *testing.T) {
	defer cleanDb(db, t)

	resRec := httptest.NewRecorder()
	body := fmt.Appendf(nil, "{\"email\": \"%s\", \"password\": \"%s\"}", mock.email, mock.pwd)
	req, err := http.NewRequest(http.MethodPost, addr+"/users/login", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	usersCtrl := getUsersController(db)
	usersCtrl.Login(resRec, req)
	if resRec.Result().StatusCode != http.StatusForbidden {
		t.Fatal(statusCodesMismatch(http.StatusForbidden, resRec.Result().StatusCode, resRec.Body.String()))
	}
}
