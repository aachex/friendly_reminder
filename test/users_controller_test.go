package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/hasher"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
	"github.com/artemwebber1/friendly_reminder/pkg/middleware"
)

var addr = cfg.Host + ":" + cfg.Port

func TestSendConfirmEmailLink(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	usersCtrl := getUsersController(db)

	body := bytes.NewReader(fmt.Appendf(nil, "{\"email\": \"%s\", \"password\": \"%s\"}", mock.email, mock.pwd))
	req, err := http.NewRequest(http.MethodPost, addr+"/new-user", body)
	if err != nil {
		t.Fatal(err)
	}

	resRec := httptest.NewRecorder()
	usersCtrl.SendConfirmEmailLink(resRec, req)

	if resRec.Result().StatusCode != 200 {
		t.Fatalf("Wanted status code 200, got %d. Body: %s", resRec.Result().StatusCode, resRec.Body)
	}

	tok := resRec.Body.String()

	uur := repository.NewUnverifiedUsersRepository(db)
	if !uur.HasToken(mock.email) || !uur.TokenExists(tok) {
		t.Fatalf("Invalid token: %s", tok)
	}
}

func TestConfirmEmail(t *testing.T) {
	db := openDb(t)
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

	uur := repository.NewUnverifiedUsersRepository(db)
	if !uur.HasToken(mock.email) || !uur.TokenExists(tok) {
		t.Fatalf("Invalid token: %s", tok)
	}

	req, err = http.NewRequest(http.MethodGet, addr+"/confirm-email?t="+tok, nil)
	if err != nil {
		t.Fatal(err)
	}
	usersCtrl.ConfirmEmail(resRec, req)

	if uur.TokenExists(tok) {
		t.Fatal("Failed to create new user")
	}
}

func TestSignUser_Unauthorized(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	usersCtrl := getUsersController(db)

	resRec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPatch, addr+"/sign-user", nil)
	if err != nil {
		t.Fatal(err)
	}

	middleware.RequireAuthorization(usersCtrl.SignUser)(resRec, req)
	if resRec.Result().StatusCode != http.StatusForbidden && resRec.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("Wanted status code %d or %d, got: %d.", http.StatusForbidden, http.StatusUnauthorized, resRec.Result().StatusCode)
	}
}

func TestLogin(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	usersRepo := repository.NewUsersRepository(db)
	_, err := usersRepo.AddUser(mock.email, hasher.Hash(mock.pwd))
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	resRec := httptest.NewRecorder()
	body := fmt.Appendf(nil, "{\"email\": \"%s\", \"password\": \"%s\"}", mock.email, mock.pwd)
	req, err := http.NewRequest(http.MethodPost, addr+"/login", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	usersCtrl := getUsersController(db)
	usersCtrl.Login(resRec, req)
	if resRec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Wanted status code %d, got %d.\nResponse body: %s", http.StatusOK, resRec.Result().StatusCode, resRec.Body)
	}
}

func TestLogin_UserDoesntExist(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	resRec := httptest.NewRecorder()
	body := fmt.Appendf(nil, "{\"email\": \"%s\", \"password\": \"%s\"}", mock.email, mock.pwd)
	req, err := http.NewRequest(http.MethodPost, addr+"/login", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	usersCtrl := getUsersController(db)
	usersCtrl.Login(resRec, req)
	if resRec.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("Wanted status code %d, got %d.\nResponse body: %s", http.StatusForbidden, resRec.Result().StatusCode, resRec.Body)
	}
}
