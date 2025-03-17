package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/hasher"
	"github.com/artemwebber1/friendly_reminder/internal/repository"
)

func TestCreateTask_Unauthorized(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	tasksCtrl := getTasksController(db)

	resRec := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPost, addr+"/new-task", nil)
	if err != nil {
		t.Fatal(err)
	}

	tasksCtrl.CreateTask(resRec, req)

	if resRec.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("Wanted status code %d, got %d", http.StatusForbidden, resRec.Result().StatusCode)
	}
}

func TestCreateTask(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	tasksCtrl := getTasksController(db)

	resRec := httptest.NewRecorder()
	body := bytes.NewReader(fmt.Appendf(nil, "{\"value\": \"%s\"}", "smth"))
	req, err := http.NewRequest(http.MethodPost, addr+"/new-task", body)
	if err != nil {
		t.Fatal(err)
	}

	usersRepo := repository.NewUsersRepository(db)
	usersRepo.AddUser(mock.email, hasher.Hash(mock.pwd))
	tok := getJwt(t, getUsersController(db))

	req.Header.Add("Authorization", "Bearer "+tok)
	tasksCtrl.CreateTask(resRec, req)

	if resRec.Result().StatusCode != http.StatusCreated {
		t.Fatal(statusCodesMismatch(http.StatusOK, resRec.Result().StatusCode, resRec.Body.String()))
	}
}

func TestDeleteTask(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	usersRepo := repository.NewUsersRepository(db)
	usersRepo.AddUser(mock.email, hasher.Hash(mock.pwd))

	tasksRepo := repository.NewTasksRepository(db)
	id, err := tasksRepo.AddTask("Do homework", mock.email)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf(addr+"/del-task?id=%d", id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	resRec := httptest.NewRecorder()

	tasksCtrl := getTasksController(db)
	tasksCtrl.DeleteTask(resRec, req)

	if resRec.Result().StatusCode != http.StatusOK {
		t.Fatal(statusCodesMismatch(http.StatusOK, resRec.Result().StatusCode, resRec.Body.String()))
	}
}
