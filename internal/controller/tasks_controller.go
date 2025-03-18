package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
	"github.com/artemwebber1/friendly_reminder/pkg/jwtservice"
	mw "github.com/artemwebber1/friendly_reminder/pkg/middleware"
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
	mux.HandleFunc(
		"POST /new-task",
		mw.UseLogging(mw.UseAuthorization(c.CreateTask)),
	)

	mux.HandleFunc(
		"GET /list",
		mw.UseLogging(mw.UseAuthorization(c.GetList)),
	)

	mux.HandleFunc(
		"DELETE /clear-list",
		mw.UseLogging(mw.UseAuthorization(c.ClearList)),
	)

	mux.HandleFunc(
		"DELETE /del-task",
		mw.UseLogging(mw.UseAuthorization(c.DeleteTask)),
	)
}

// CreateTask создаёт новую задачу в списке пользователя.
//
// Обрабатывает POST запросы по пути '/new-task'.
func (c *TasksController) CreateTask(w http.ResponseWriter, r *http.Request) {
	rawJwt := getRawJwtFromHeader(r.Header)
	jwtClaims, err := jwtservice.GetClaims(rawJwt, jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	email, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	type newTask struct {
		Id    int64  `json:"task_id"`
		Value string `json:"value"`
	}

	task, err := readBody[newTask](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*30)
	defer cancel()

	id, err := c.tasksRepo.AddTask(ctx, task.Value, email)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == context.DeadlineExceeded {
			// Контекст отменён по таймауту
			statusCode = http.StatusRequestTimeout
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	task.Id = id

	w.WriteHeader(http.StatusCreated)
	writeJson(w, task)
}

func (c *TasksController) GetList(w http.ResponseWriter, r *http.Request) {
	rawJwt := getRawJwtFromHeader(r.Header)
	jwtClaims, err := jwtservice.GetClaims(rawJwt, jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	email, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*50)
	defer cancel()

	list, err := c.tasksRepo.GetList(ctx, email)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == context.DeadlineExceeded {
			// Контекст отменён по таймауту
			statusCode = http.StatusRequestTimeout
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	writeJson(w, &list)
}

func (c *TasksController) ClearList(w http.ResponseWriter, r *http.Request) {
	rawJwt := getRawJwtFromHeader(r.Header)
	jwtClaims, err := jwtservice.GetClaims(rawJwt, jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	email, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*20)
	defer cancel()

	err = c.tasksRepo.ClearList(ctx, email)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == context.DeadlineExceeded {
			// Контекст отменён по таймауту
			statusCode = http.StatusRequestTimeout
		}
		http.Error(w, err.Error(), statusCode)
		return
	}
}

// DeleteTask удаляет задачу из списка пользователя.
//
// Обрабатывает DELETE запросы по пути '/del-task'.
func (c *TasksController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*20)
	defer cancel()

	err = c.tasksRepo.DeleteTask(ctx, taskId)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == context.DeadlineExceeded {
			// Контекст отменён по таймауту
			statusCode = http.StatusRequestTimeout
		}
		http.Error(w, err.Error(), statusCode)
		return
	}
}
