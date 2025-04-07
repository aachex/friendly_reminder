package test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/hasher"
	repo "github.com/artemwebber1/friendly_reminder/internal/repository/sqlite"
)

func TestCreateTask_Unauthorized(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	tasksCtrl := getTasksController(db)

	resRec := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPost, addr+"/tasks/new", nil)
	if err != nil {
		t.Fatal(err)
	}

	tasksCtrl.CreateTask(resRec, req)

	if resRec.Result().StatusCode != http.StatusForbidden {
		t.Fatal(statusCodesMismatch(http.StatusForbidden, resRec.Result().StatusCode, resRec.Body.String()))
	}
}

func TestCreateTask(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	tasksCtrl := getTasksController(db)

	resRec := httptest.NewRecorder()
	body := bytes.NewReader(fmt.Appendf(nil, "{\"value\": \"%s\"}", "smth"))
	req, err := http.NewRequest(http.MethodPost, addr+"/tasks/new", body)
	if err != nil {
		t.Fatal(err)
	}

	usersRepo := repo.NewUsersRepository(db)
	usersRepo.AddUser(t.Context(), mock.email, hasher.Hash(mock.pwd))

	req.Header.Add("Authorization", "Bearer "+getJwt(t, getUsersController(db)))
	tasksCtrl.CreateTask(resRec, req)

	if resRec.Result().StatusCode != http.StatusCreated {
		t.Fatal(statusCodesMismatch(http.StatusOK, resRec.Result().StatusCode, resRec.Body.String()))
	}
}

func TestDeleteTask(t *testing.T) {
	db := openDb(t)
	defer cleanDb(db, t)

	usersRepo := repo.NewUsersRepository(db)
	usersRepo.AddUser(t.Context(), mock.email, hasher.Hash(mock.pwd))

	tasksRepo := repo.NewTasksRepository(db)

	ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond*20)
	defer cancel()

	id, err := tasksRepo.AddTask(ctx, "Do homework", mock.email)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf(addr+"/tasks/del?id=%d", id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	resRec := httptest.NewRecorder()

	req.Header.Add("Authorization", "Bearer "+getJwt(t, getUsersController(db)))
	tasksCtrl := getTasksController(db)
	tasksCtrl.DeleteTask(resRec, req)

	if resRec.Result().StatusCode != http.StatusOK {
		t.Fatal(statusCodesMismatch(http.StatusOK, resRec.Result().StatusCode, resRec.Body.String()))
	}
}
