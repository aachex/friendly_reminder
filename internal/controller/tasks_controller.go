package controller

import (
	"net/http"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
	"github.com/artemwebber1/friendly_reminder/pkg/middleware"
)

type TasksController struct {
	tasksRepo repository.TasksRepository
}

func NewTasksController(tr repository.TasksRepository) *TasksController {
	return &TasksController{
		tasksRepo: tr,
	}
}

func (c *TasksController) AddEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/new-task", middleware.RequireAuthorization(c.CreateTask))
}

func (c *TasksController) CreateTask(w http.ResponseWriter, r *http.Request) {
	rawTok := getRawJwtFromHeader(r.Header)
	jwtClaims, err := readJWT(rawTok)
	if err != nil {
		http.Error(w, errInvalidToken.Error(), http.StatusForbidden)
		return
	}

	email, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, errInvalidTokenSubject.Error(), http.StatusForbidden)
		return
	}

	type newTask struct {
		Value string `json:"value"`
	}

	task, err := readBody[newTask](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.tasksRepo.AddTask(email, task.Value)
	w.WriteHeader(http.StatusCreated)
}
