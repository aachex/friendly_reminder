package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/artemwebber1/friendly_reminder/internal/hasher"
	repo "github.com/artemwebber1/friendly_reminder/internal/repository/postgres"
)

func TestCreateTask_Unauthorized(t *testing.T) {
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
	defer cleanDb(db, t)

	tasksCtrl := getTasksController(db)

	resRec := httptest.NewRecorder()
	body := bytes.NewReader(fmt.Appendf(nil, "{\"value\": \"%s\"}", "smth"))
	req, err := http.NewRequest(http.MethodPost, addr+"/tasks/new", body)
	if err != nil {
		t.Fatal(err)
	}

	usersRepo := repo.NewUsersRepository(db)
	err = usersRepo.AddUser(t.Context(), mock.email, hasher.Hash(mock.pwd))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+getJwt(t, getUsersController(db)))
	tasksCtrl.CreateTask(resRec, req)

	if resRec.Result().StatusCode != http.StatusCreated {
		t.Fatal(statusCodesMismatch(http.StatusOK, resRec.Result().StatusCode, resRec.Body.String()))
	}
}

func TestDeleteTask(t *testing.T) {
	defer cleanDb(db, t)

	usersRepo := repo.NewUsersRepository(db)
	usersRepo.AddUser(t.Context(), mock.email, hasher.Hash(mock.pwd))

	tasksRepo := repo.NewTasksRepository(db)

	id, err := tasksRepo.AddTask(t.Context(), "Do homework", mock.email)
	if err != nil {
		t.Fatal(err)
	}

	url := addr + "/tasks/del/{id}"
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetPathValue("id", strconv.FormatInt(id, 10))

	resRec := httptest.NewRecorder()

	req.Header.Add("Authorization", "Bearer "+getJwt(t, getUsersController(db)))
	tasksCtrl := getTasksController(db)
	tasksCtrl.DeleteTask(resRec, req)

	if resRec.Result().StatusCode != http.StatusOK {
		t.Fatal(statusCodesMismatch(http.StatusOK, resRec.Result().StatusCode, resRec.Body.String()))
	}
}
