package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

func TestSendConfirmEmailLink(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	usersCtrl := getUsersController(db)

	url := cfg.Host + ":" + cfg.Port
	body := bytes.NewReader(fmt.Appendf(nil, "{\"email\": \"%s\", \"password\": \"%s\"}", mock.email, mock.pwd))
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		t.Fatal(err)
	}

	resRec := httptest.NewRecorder()
	usersCtrl.SendConfirmEmailLink(resRec, req)

	if resRec.Code != 200 {
		t.Fatalf("Wanted status code 200, got %d. Body: %s", resRec.Code, resRec.Body)
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

	addr := cfg.Host + ":" + cfg.Port
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
